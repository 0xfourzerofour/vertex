### Vertex GraphQL Api Gateway ### 

This Project is inspired by current pitfalls that I have come across at work with monolithic graphql schemas.
Vertex aims to solve this issue by allowing a single graphql endpoint to many downstream services by parsing 
the query body and matching the query to a service. 

Vertex is designed to be a serverless Graphql api gateway that can run on the AWS Lambda platform. All the proxy logic is called
by the vertex handler. All the handler needs is the context to a vertex config (service query map, introspection schema, fasthttp client)

[![](https://mermaid.ink/img/pako:eNp9j00KwjAQRq8SZt1eIAtB2qgLF9UGFZouQhNtsElKTMHS9u6m_iwEcRbD8L3HMDNAZYUEDBfH2xrRlBkUalkkjZLGlyiOF-OG0iyaWx4d8xElw4HsKTlNLzd5Opmz935EabGeF-225Q9I_sHVF4QItHSaKxFuG-aEga-llgxwGAV3VwbMTMHrWsG9JEJ56wCfeXOTEfDO27w3FWDvOvmRUsXDn_ptTQ8RuVGX)](https://mermaid-js.github.io/mermaid-live-editor/edit#pako:eNp9j00KwjAQRq8SZt1eIAtB2qgLF9UGFZouQhNtsElKTMHS9u6m_iwEcRbD8L3HMDNAZYUEDBfH2xrRlBkUalkkjZLGlyiOF-OG0iyaWx4d8xElw4HsKTlNLzd5Opmz935EabGeF-225Q9I_sHVF4QItHSaKxFuG-aEga-llgxwGAV3VwbMTMHrWsG9JEJ56wCfeXOTEfDO27w3FWDvOvmRUsXDn_ptTQ8RuVGX)

# NOTE: THIS IS STILL VERY EXPERIMENTAL SO DO NOT USE IN PRODUCTION ANYTIME SOON

## Example Lambda Proxy Integration

```
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
		"allFilms":  "swapi-graphql.netlify.app/.netlify/functions/index",
	}

	vert := vertex.NewVertex(exampleMap, schema, clients.FHTTP())

	lambda.Start(vert.Handler)
}

```

# Vertex CLI

A CLI tool that generates a JSON file to be used with the Vertex Graphql Proxy

The `generate` command accepts 4 flags

 - --schema || -s
 - --output || -o
 - --name || -n
 - --url || -u

It will generate the following JSON file

```
{
    "serviceName": "example-service-name",
    "serviceUrl": "service.com/graphql",
    "schema": "*** graphql schema definition",
    "queryMap": {
      "exampleQuery": "service.com/graphql"
    }
}

```

## Exciting findings

 - requests can be parsed within 0.01-0.1ms which will only get faster after I implement request caching

## Considerations

 - There cannot be overlapping types throughout the graphql schemas (Might be an option for query polymorphism based on different variables)
 - Services must have introspection turned on at the API level (I am working around this using a ci tool to publish schemas)
 - This project is still early days so do not use in production

## TODO

 - Create cli tool to publish schemas from service CI [current tool only for generation](https://github.com/joshpauline/vertex-cli) 
 - create cli tool to deploy schemas from service ci to remove need for introspection
 - Load all graphs into one schema for playground introspection
 - Add Web socket functionality for subsriptions
 - cache full introspection graph in memory for fast retrieval
 - Cache hash of request body for faster proxying
 - Load and save services from DynamoDB
 - Logging
 - Load services from CLI
