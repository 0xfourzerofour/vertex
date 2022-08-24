package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	log.Print("LOAD SERVICES INTO CACHE")

}

func Handle(ctx context.Context, input interface{}) (*string, error) {

	return nil, nil
}

func main() {
	lambda.Start(Handle)
}
