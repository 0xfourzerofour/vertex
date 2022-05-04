package main

import (
	"encoding/json"
	"govertex/internal/graphql"
	"govertex/internal/service"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/valyala/fasthttp"
)

type GQLResp struct {
	Data  map[string]interface{} `json:"data"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func ProxyHandler(ctx *fasthttp.RequestCtx) {

	body := ctx.Request.Body()

	queries, err := graphql.ParseQueryBody(&body)

	if err != nil {
		log.Print("Could not parse request")
	}

	resChannel := make(chan GQLResp)
	errChanel := make(chan []byte)
	d := make(map[string]interface{})
	result := GQLResp{
		Data: d,
	}

	for _, query := range *queries {

		if proxy, ok := service.ProxyMap.GetStringKey(query); ok {

			proxyStr := proxy.(string)

			if path, ok := service.ServiceMap.GetStringKey(query); ok {
				proxyStr += path.(string)
			}

			go func() {

				sendPostRequest(service.Client, proxyStr, &body, resChannel, errChanel)

			}()
		}
	}

	for _, query := range *queries {

		data := <-resChannel

		result.Data[query] = data.Data[query]

	}

	final, err := json.Marshal(result)

	ctx.Response.SetBody(final)

}

func main() {

	log.Print("Starting")

	err := service.LoadServices()

	if err != nil {
		log.Fatal("Could not load services")
	}

	if err := fasthttp.ListenAndServe("localhost:3000", ProxyHandler); err != nil {
		log.Fatal(err)
	}

}

var headerContentTypeJson = []byte("application/json")

var client *fasthttp.Client

type Entity struct {
	Id   int
	Name string
}

func sendPostRequest(client *fasthttp.Client, url string, body *[]byte, bodyBytes chan<- GQLResp, errorChan chan<- []byte) {
	// per-request timeout
	reqTimeout := 5 * time.Second

	log.Print(url)

	req := fasthttp.AcquireRequest()
	req.SetRequestURI("https://" + url)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentTypeBytes(headerContentTypeJson)
	req.SetBodyRaw(*body)
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

	resBody := GQLResp{}

	err = json.Unmarshal(resp.Body(), &resBody)

	bodyBytes <- resBody

	fasthttp.ReleaseResponse(resp)
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