package usecase

import (
	"context"
	"encoding/json"
	"govertex/domain/proxy"
	"govertex/domain/schemas"
	"log"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

type vertexUsecase struct {
	schemaRepo schemas.SchemaRepository
	proxyRepo  proxy.ProxyRepository
}

func NewVertexUsecase(schemaRepo schemas.SchemaRepository, proxyRepo proxy.ProxyRepository) *vertexUsecase {
	return &vertexUsecase{
		schemaRepo,
		proxyRepo,
	}
}

func (v *vertexUsecase) MergeSchemas(ctx context.Context) error {

	schemaList, err := v.schemaRepo.ListSubSchemas(ctx)

	if err != nil {
		return err
	}

	return v.schemaRepo.Merge(ctx, schemaList)
}

func (v *vertexUsecase) ProxyHandler(fastctx *fasthttp.RequestCtx) {
	body := fastctx.Request.Body()

	queries, err := schemas.ParseQueryBody(&body)

	if err != nil {
		log.Print("Could not parse request")
	}

	result, err := v.proxyRepo.SendConcurrentRequests(fastctx, queries)

	final, err := json.Marshal(result)

	if err != nil {
		fastctx.Response.SetStatusCode(500)
	}

	fastctx.Response.SetBody(final)
}

func (v *vertexUsecase) Listen(ctx context.Context) {

	inmemlistener := fasthttputil.NewInmemoryListener()

	if err := fasthttp.Serve(inmemlistener, v.ProxyHandler); err != nil {
		log.Fatal(err)
	}

}
