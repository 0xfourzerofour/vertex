package graphql

import (
	"encoding/json"
	"log"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
)

type HttpQuery struct {
	Query string `json:"query"`
}

func GetQueryService(body *[]byte) {

	parseQuery(body)

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
