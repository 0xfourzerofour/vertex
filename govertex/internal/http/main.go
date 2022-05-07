package http

import (
	"context"
	"encoding/json"
	"govertex/internal/graphql"
	"govertex/internal/service"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/valyala/fasthttp"
	"golang.org/x/sync/errgroup"
)

type GQLResp struct {
	Data  map[string]interface{} `json:"data,omitempty"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func SendConcurrentQueries(queries []*graphql.SubQuery) (*GQLResp, error) {

	ctx := context.Background()

	g, _ := errgroup.WithContext(ctx)

	d := make(map[string]interface{})

	result := GQLResp{
		Data: d,
	}

	for _, query := range queries {

		queryName := query.QueryName

		if proxy, ok := service.ProxyMap.GetStringKey(query.QueryName); ok {

			proxyStr := proxy.(string)

			if path, ok := service.ServiceMap.GetStringKey(query.QueryName); ok {
				proxyStr += path.(string)
			}

			b, err := json.Marshal(query.Body)

			if err != nil {
				log.Print(err)
			}

			g.Go(func() error {

				postRes, err := sendPostRequest(service.Client, proxyStr, b)

				if err == nil {
					result.Data[queryName] = postRes.Data[queryName]
				}

				return err

			})
		}
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &result, nil

}

var headerContentTypeJson = []byte("application/json")

var client *fasthttp.Client

func sendPostRequest(client *fasthttp.Client, url string, body []byte) (*GQLResp, error) {

	reqTimeout := 5 * time.Second

	log.Print("URL: ", url, string(body))

	req := fasthttp.AcquireRequest()
	req.SetRequestURI("https://" + url)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentTypeBytes(headerContentTypeJson)
	req.SetBodyRaw(body)
	resp := fasthttp.AcquireResponse()
	err := client.DoTimeout(req, resp, reqTimeout)
	fasthttp.ReleaseRequest(req)
	if err == nil {
		statusCode := resp.StatusCode()
		if statusCode == http.StatusOK {
			log.Print("SUCCESS")
		} else {
			log.Printf("ERR invalid HTTP response code: %d\n", statusCode)
		}
	} else {

		errName, known := httpConnError(err)
		log.Print(errName)
		if known {
			log.Print(known)
		} else {
			log.Print("Unkownerror")
		}
	}

	bodyVal := resp.Body()

	fasthttp.ReleaseResponse(resp)

	resBody := GQLResp{}

	err = json.Unmarshal(bodyVal, &resBody)

	if err != nil {
		return nil, err

	}

	return &resBody, nil

}

func httpConnError(err error) (string, bool) {
	errName := ""
	known := false
	if err == fasthttp.ErrTimeout {
		errName = "timeout"
		known = true
	} else if err == fasthttp.ErrNoFreeConns {
		errName = "conn_limit"
		known = true
	} else if err == fasthttp.ErrConnectionClosed {
		errName = "conn_close"
		known = true
	} else {
		errName = reflect.TypeOf(err).String()
		if errName == "*net.OpError" {
			// Write and Read errors are not so often and in fact they just mean timeout problems
			errName = "timeout"
			known = true
		}
	}
	return errName, known
}