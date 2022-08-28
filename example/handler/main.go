package main

import (
	_ "embed"
	"github.com/joshpauline/vertex/internal/clients"
	"github.com/joshpauline/vertex/pkg/vertex"

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
