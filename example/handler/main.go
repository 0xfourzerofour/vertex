package main

import (
	_ "embed"
	"govertex/internal/clients"
	"govertex/pkg/vertex"

	"github.com/aws/aws-lambda-go/lambda"
)

//go:embed schema.graphql
var schema string

func main() {

	exampleMap := map[string]string{
		"countries": "countries.trevorblades.com/graphql",
	}

	vert := vertex.NewVertex(exampleMap, schema, clients.FHTTP())

	lambda.Start(vert.Handler)
}
