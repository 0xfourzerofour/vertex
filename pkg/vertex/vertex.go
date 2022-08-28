package vertex

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/graph-gophers/graphql-go"
	"github.com/valyala/fasthttp"
)

type vertex struct {
	vertexMap map[string]string
	schema    *graphql.Schema
	client    *fasthttp.Client
}

func NewVertex(vertexMap map[string]string, schemaStr string, client *fasthttp.Client) *vertex {
	schema := graphql.MustParseSchema(schemaStr, nil)
	return &vertex{
		vertexMap: vertexMap,
		schema:    schema,
		client:    client,
	}
}

func (v *vertex) Handler(ctx context.Context, input events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	queries, err := v.parseQueryBody(input.Body)

	if err != nil {
		return nil, err
	}

	res, err := v.sendConcurrentRequests(ctx, input, queries)

	if err != nil {
		return nil, err
	}

	jsonResult, err := json.Marshal(res)

	if err != nil {
		return nil, err
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonResult),
	}, nil
}
