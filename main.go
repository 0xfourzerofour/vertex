package main

import (
	"govertex/internal/graphql"
	"govertex/internal/service"
	"log"
	"time"

	"github.com/valyala/fasthttp"

	proxy "github.com/yeqown/fasthttp-reverse-proxy/v2"
)

func ProxyHandler(ctx *fasthttp.RequestCtx) {

	start := time.Now()

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

		t := time.Now()

		log.Println(t.Sub(start))

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