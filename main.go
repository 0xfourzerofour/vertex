package main

import (
	"govertex/internal/graphql"
	"govertex/internal/service"
	"log"

	"github.com/valyala/fasthttp"

	proxy "github.com/yeqown/fasthttp-reverse-proxy/v2"
)

var (
	proxyServer  = proxy.NewReverseProxy("localhost:8000")
	proxyServer1 = proxy.NewReverseProxy("localhost:3500")
	proxyServer2 = proxy.NewReverseProxy("localhost:3001")
)

func ProxyHandler(ctx *fasthttp.RequestCtx) {

	body := ctx.Request.Body()

	graphql.ParseQueryBody(&body)

	proxyServer.ServeHTTP(ctx)

}

func main() {
	// Load services from config

	err := service.LoadServices()

	if err != nil {
		log.Fatal("Could not load services")
	}

	// build proxies

	//introspection must be on to get schema docs

	if err := fasthttp.ListenAndServe("localhost:3000", ProxyHandler); err != nil {
		log.Fatal(err)
	}
}