package main

import (
	"context"
	"govertex/application/usecase"
	"govertex/infra/persistance"
	"govertex/internal/clients"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	log.Print("LOAD SERVICES INTO CACHE")

}

func Handle(ctx context.Context, input interface{}) {

	pers := persistance.ProxyPersistance(clients.FHTTP())
	schemaPers := persistance.SchemaPersistance(clients.S3())
	usecase := usecase.NewVertexUsecase(schemaPers, pers)

	usecase.Listen(ctx)

}

func main() {
	lambda.Start(Handle)
}
