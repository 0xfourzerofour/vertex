package vertex

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/valyala/fasthttp"

	"golang.org/x/sync/errgroup"
)

type GQLResp struct {
	Data  map[string]interface{} `json:"data,omitempty"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (v *vertex) httpConnError(err error) (string, bool) {
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
			errName = "timeout"
			known = true
		}
	}
	return errName, known
}

func (v *vertex) forward(ctx context.Context, originReq events.APIGatewayProxyRequest, url string, body []byte) (*GQLResp, error) {

	reqTimeout := 5 * time.Second

	req := fasthttp.AcquireRequest()

	var headerContentTypeJson = []byte("application/json")

	req.SetRequestURI("https://" + url)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentTypeBytes(headerContentTypeJson)

	for key, header := range originReq.Headers {
		req.Header.Add(key, header)
	}

	req.SetBodyRaw(body)
	resp := fasthttp.AcquireResponse()

	err := v.client.DoTimeout(req, resp, reqTimeout)
	fasthttp.ReleaseRequest(req)

	if err == nil {
		statusCode := resp.StatusCode()
		if statusCode == http.StatusOK {
			log.Print("SUCCESS")
		} else {
			log.Printf("ERR invalid HTTP response code: %d\n", statusCode)
		}
	} else {

		errName, known := v.httpConnError(err)
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

func (v *vertex) sendConcurrentRequests(ctx context.Context, originReq events.APIGatewayProxyRequest, queries []*subQuery) (*GQLResp, error) {

	g := errgroup.Group{}

	d := make(map[string]interface{})

	result := GQLResp{
		Data: d,
	}

	for _, query := range queries {

		queryName := query.QueryName

		if proxy, ok := v.vertexMap[query.QueryName]; ok {

			b, err := json.Marshal(query.Body)

			if err != nil {
				log.Print(err)
			}

			g.Go(func() error {

				postRes, err := v.forward(ctx, originReq, proxy, b)

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
