package main

import (
	"context"
	"govertex/application/usecase"
	"govertex/infra/persistance"
	"govertex/internal/clients"

	"github.com/aws/aws-lambda-go/lambda"
)

func Handle(ctx context.Context, input interface{}) error {
	pers := persistance.SchemaPersistance(clients.S3())
	usecase := usecase.NewVertexUsecase(pers, nil)
	return usecase.MergeSchemas(ctx)
}

func main() {
	lambda.Start(Handle)
}
