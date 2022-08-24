package usecase

import (
	"context"
	"govertex/domain/proxy"
	"govertex/domain/schemas"
	"log"

	"golang.org/x/sync/errgroup"
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

	log.Print("SCHEMAS: ", schemaList)

	return v.schemaRepo.Merge(ctx, schemaList)
}

func (v *vertexUsecase) SendConcurrentRequests(ctx context.Context, queries []*graphql.SubQuery) error {

	g, _ := errgroup.WithContext(ctx)

	d := make(map[string]interface{})

	result := proxy.GQLResp{
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
