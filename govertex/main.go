package main

import (
	"encoding/json"
	"govertex/internal/graphql"
	internalHttp "govertex/internal/http"

	"govertex/internal/service"
	"log"
	"time"

	"github.com/valyala/fasthttp"
)

func ProxyHandler(fastctx *fasthttp.RequestCtx) {

	init := time.Now()

	body := fastctx.Request.Body()

	queries, err := graphql.ParseQueryBody(&body)

	if err != nil {
		log.Print("Could not parse request")
	}

	after := time.Now()

	result, err := internalHttp.SendConcurrentQueries(queries)

	nanoTime := float64(after.Nanosecond()) - float64(init.Nanosecond())
	log.Printf("%d QUERIES PARSED in %f ns OR %f ms", len(queries), nanoTime, nanoTime/1000000)

	final, err := json.Marshal(result)

	if err != nil {
		fastctx.Response.SetStatusCode(500)
	}

	fastctx.Response.SetBody(final)

}

func main() {

	err := service.LoadServices()

	if err != nil {
		log.Fatal("Could not load services")
	}

	log.Print("SERVICES LOADED")

	if err := fasthttp.ListenAndServe("localhost:3000", ProxyHandler); err != nil {
		log.Fatal(err)
	}

}
