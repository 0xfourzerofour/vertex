package usecase

import (
	"context"
	"govertex/domain/proxy"
	"govertex/domain/schemas"
	"log"
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
