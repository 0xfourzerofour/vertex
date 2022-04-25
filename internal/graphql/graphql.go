package graphql

import (
	"context"
	"encoding/json"
	"log"

	"github.com/machinebox/graphql"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
)

type HttpQuery struct {
	Query string `json:"query"`
}

type IntroSpectionResult struct {
	Data SchemaResult `json:"data"`
}
type SchemaResult struct {
	Schema SchemaTypes `json:"__schema"`
}

type SchemaTypes struct {
	Types []TypeData `json:"types"`
}

type TypeData struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

func GetQueryService(body *[]byte) {

	parseQuery(body)

}
func GetIntrospectionSchema(url string) (*IntroSpectionResult, error) {

	req := graphql.NewRequest(`
		{
		  __schema {
		    types {
		      name
		      description
		    }
		  }
		}		
	`)

	client := graphql.NewClient(url)

	resp := IntroSpectionResult{}

	err := client.Run(context.Background(), req, &resp)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func parseQuery(body *[]byte) {

	hq := HttpQuery{}

	err := json.Unmarshal(*body, &hq)

	if err != nil {
		log.Print(err)
	}

	queryAst := ast.Source{
		Input: hq.Query,
	}

	as, err := parser.ParseQuery(&queryAst)

	if err != nil {
		log.Print(err)
	}

	for _, operation := range as.Operations {
		// Name query/mutation name

		//Directive query or mutation

		log.Printf("%+v", operation.SelectionSet[0])
	}

}
