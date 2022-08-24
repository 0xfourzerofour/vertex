package persistance

import (
	"context"
	"encoding/json"
	"govertex/domain/proxy"
	"govertex/domain/schemas"
	"log"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
	"golang.org/x/sync/errgroup"
)

type proxyImp struct {
	proxyConn *fasthttp.Client
}

func ProxyPersistance(proxyConn *fasthttp.Client) proxy.ProxyRepository {
	return &proxyImp{
		proxyConn,
	}
}

func (p *proxyImp) Forward(ctx context.Context, url string, body []byte) (*proxy.GQLResp, error) {

	reqTimeout := 5 * time.Second

	req := fasthttp.AcquireRequest()

	var headerContentTypeJson = []byte("application/json")

	req.SetRequestURI("https://" + url)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentTypeBytes(headerContentTypeJson)
	req.SetBodyRaw(body)
	resp := fasthttp.AcquireResponse()

	err := p.proxyConn.DoTimeout(req, resp, reqTimeout)
	fasthttp.ReleaseRequest(req)

	if err == nil {
		statusCode := resp.StatusCode()
		if statusCode == http.StatusOK {
			log.Print("SUCCESS")
		} else {
			log.Printf("ERR invalid HTTP response code: %d\n", statusCode)
		}
	} else {

		errName, known := proxy.HttpConnError(err)
		log.Print(errName)
		if known {
			log.Print(known)
		} else {
			log.Print("Unkownerror")
		}
	}

	bodyVal := resp.Body()

	fasthttp.ReleaseResponse(resp)

	resBody := proxy.GQLResp{}

	err = json.Unmarshal(bodyVal, &resBody)

	if err != nil {
		return nil, err
	}

	return &resBody, nil
}

func (p *proxyImp) SendConcurrentRequests(ctx *fasthttp.RequestCtx, queries []*schemas.SubQuery) (*proxy.GQLResp, error) {

	g := errgroup.Group{}

	d := make(map[string]interface{})

	result := proxy.GQLResp{
		Data: d,
	}

	for _, query := range queries {

		queryName := query.QueryName

		log.Print(queryName)

		// if proxy, ok := service.ProxyMap.GetStringKey(query.QueryName); ok {

		// 	proxyStr := proxy.(string)

		// 	// if path, ok := service.ServiceMap.GetStringKey(query.QueryName); ok {

		// 	// 	proxyStr += path.(string)

		// 	// }

		// 	b, err := json.Marshal(query.Body)

		// 	if err != nil {
		// 		log.Print(err)
		// 	}

		// 	g.Go(func() error {

		// 		postRes, err := p.Forward(ctx, proxyStr, b)

		// 		if err == nil {
		// 			result.Data[queryName] = postRes.Data[queryName]
		// 		}

		// 		return err

		// 	})
		// }
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &result, nil
}
