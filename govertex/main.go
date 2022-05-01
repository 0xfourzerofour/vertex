package main

import (
	"github.com/valyala/fasthttp"
	"govertex/internal/graphql"
	"govertex/internal/service"
	"log"

	proxy "github.com/yeqown/fasthttp-reverse-proxy/v2"
)

func ProxyHandler(ctx *fasthttp.RequestCtx) {

	body := ctx.Request.Body()

	query, err := graphql.ParseQueryBody(&body)

	if err != nil {
		log.Print("Could not parse request")
	}

	proxyServer, _ := service.ProxyMap.GetStringKey(*query)

	if proxyServer != nil {

		if field, ok := service.ServiceMap.GetStringKey(*query); ok {

			ctx.Request.SetRequestURI(field.(string))
		}

		proxyServer.(*proxy.ReverseProxy).ServeHTTP(ctx)

		return
	}

	log.Print("Could not proxy request")

}

func main() {

	err := service.LoadServices()

	if err != nil {
		log.Fatal("Could not load services")
	}

	if err := fasthttp.ListenAndServe("localhost:3000", ProxyHandler); err != nil {
		log.Fatal(err)
	}

}