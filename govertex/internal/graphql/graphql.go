package graphql

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"golang.org/x/oauth2"

	"github.com/shurcooL/graphql"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
)

type HttpQuery struct {
	Query string `json:"query"`
}

type IntroSpectionResult struct {
	Schema struct {
		QueryType struct {
			Fields []struct {
				Name        string  `graphql:"name"`
				Description *string `graphql:"description"`
			}
		} `graphql:"queryType"`
		MutationType struct {
			Fields []struct {
				Name        string  `graphql:"name"`
				Description *string `graphql:"description"`
			}
		} `graphql:"mutationType"`
		SubscriptionType struct {
			Fields []struct {
				Name        string  `graphql:"name"`
				Description *string `graphql:"description"`
			}
		} `graphql:"subscriptionType"`
	} `graphql:"__schema"`
}

func GetIntrospectionSchema(url, token string) (*IntroSpectionResult, error) {

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := graphql.NewClient(url, httpClient)

	query := IntroSpectionResult{}

	err := client.Query(context.Background(), &query, nil)

	if err != nil {
		return nil, err
	}

	return &query, nil
}

func ParseQueryBody(body *[]byte) (*[]string, error) {

	hq := HttpQuery{}

	err := json.Unmarshal(*body, &hq)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	queryAst := ast.Source{
		Input: hq.Query,
	}

	qAst, err := parser.ParseQuery(&queryAst)

	if err.Error() != "" {
		return nil, err
	}

	queries := []string{}

	for _, operation := range qAst.Operations {

		log.Print("OPERATION: ", operation.Operation)

		log.Print(operation.SelectionSet)

		for i, selection := range operation.SelectionSet {

			field := selection.(*ast.Field)
			queries = append(queries, field.Name)

			if i == len(operation.SelectionSet)-1 {
				lastQuery := hq.Query[field.Position.Start:]

				lastBrace := strings.LastIndex(lastQuery, "}")

				lastQuery = lastQuery[:lastBrace]

				log.Print(lastQuery)

				break
			}

			if i > 0 {

				prevField := operation.SelectionSet[i-1].(*ast.Field)

				previosQuery := hq.Query[prevField.Position.Start:field.Position.Start]

				log.Print(previosQuery)

			}

		}

	}

	return &queries, nil

}
