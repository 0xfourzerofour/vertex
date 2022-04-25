package main

import (
	"encoding/json"
	"log"

	"github.com/valyala/fasthttp"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"

	proxy "github.com/yeqown/fasthttp-reverse-proxy/v2"
)

var (
	proxyServer  = proxy.NewReverseProxy("localhost:8000")
	proxyServer1 = proxy.NewReverseProxy("localhost:3500")
	proxyServer2 = proxy.NewReverseProxy("localhost:3001")
)

type HttpQuery struct {
	Query string `json:"query"`
}

// ProxyHandler ... fasthttp.RequestHandler func
func ProxyHandler(ctx *fasthttp.RequestCtx) {
	// parse gql and send to correct service

	body := ctx.Request.Body()

	hq := HttpQuery{}

	err := json.Unmarshal(body, &hq)

	if err != nil {
		log.Print(err)
	}

	queryAst := ast.Source{
		Input: hq.Query,
	}

	as, err := parser.ParseQuery(&queryAst)

	if err != nil {
		log.Print(err)
	}

	for _, operation := range as.Operations {
		// Name query/mutation name

		//Directive query or mutation

		log.Printf("%+v", operation.SelectionSet[0])
	}

	proxyServer.ServeHTTP(ctx)

}

func main() {
	// Load services from config

	// build proxies

	if err := fasthttp.ListenAndServe("localhost:3000", ProxyHandler); err != nil {
		log.Fatal(err)
	}
}