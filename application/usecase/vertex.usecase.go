package usecase

import (
	"context"
	"encoding/json"
	"govertex/domain/proxy"
	"govertex/domain/schemas"
	"log"

	"github.com/aws/aws-lambda-go/events"
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

func (v *vertexUsecase) ProxyHandler(ctx context.Context, req events.APIGatewayProxyRequest) ([]byte, error) {
	body := []byte(req.Body)

	queries, err := schemas.ParseQueryBody(&body)

	if err != nil {
		log.Print("Could not parse request")
	}

	result, err := v.proxyRepo.SendConcurrentRequests(ctx, req, queries)

	final, err := json.Marshal(result)

	if err != nil {
		return nil, err

	}

	return final, nil

}
