package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context) (Response, error) {

	resp := Response{
		StatusCode: 301,
		Headers: map[string]string{
			"Location": "countries.trevorblades.com/graphql",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
