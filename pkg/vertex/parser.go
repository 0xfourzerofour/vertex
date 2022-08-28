package vertex

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
)

type httpQuery struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type variableDirectives struct {
	Name       string
	IsNullable bool
}

type subQuery struct {
	QueryName string `json:"queryName"`
	Query     string `json:"query"`
	Operation string `json:"operation"`
	Body      httpQuery
}

type VertexData struct {
	ServiceName string            `json:"serviceName"`
	ServiceUrl  string            `json:"serviceUrl"`
	Schema      string            `json:"schema"`
	QueryMap    map[string]string `json:"queryMap"`
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

func (v *vertex) parseQueryBody(body string) ([]*subQuery, error) {

	hq := httpQuery{}

	err := json.Unmarshal([]byte(body), &hq)

	if err != nil {
		return nil, err
	}

	queryAst := ast.Source{
		Input: hq.Query,
	}

	qAst, err := parser.ParseQuery(&queryAst)

	if err.Error() != "" {
		return nil, err
	}

	queries := []*subQuery{}

	for _, operation := range qAst.Operations {

		operationMap := make(map[string]*variableDirectives)

		for _, variableDef := range operation.VariableDefinitions {
			operationMap[variableDef.Variable] = &variableDirectives{
				Name:       variableDef.Type.Name(),
				IsNullable: variableDef.Type.NonNull,
			}

		}

		for i, selection := range operation.SelectionSet {

			field := selection.(*ast.Field)

			if i == len(operation.SelectionSet)-1 {
				lastQuery := hq.Query[field.Position.Start:]

				lastBrace := strings.LastIndex(lastQuery, "}")

				lastQuery = lastQuery[:lastBrace]

				newBody := fmt.Sprintf(`%s %s { %s }`, operation.Operation, operation.Name, lastQuery)

				variables := make(map[string]interface{})

				variableStr := ""

				for i, queryVar := range field.Arguments {
					if val, ok := hq.Variables[queryVar.Name]; ok {
						variables[queryVar.Name] = val
					}
					if val, ok := operationMap[queryVar.Name]; ok {

						variableStr += "$" + queryVar.Name + ": " + val.Name
						if val.IsNullable {
							variableStr += "!"
						}

						if i != len(field.Arguments)-1 {
							variableStr += ","
						}
					}
				}

				if variableStr != "" {
					newBody = fmt.Sprintf(`%s %s(%s) { %s }`, operation.Operation, operation.Name, variableStr, lastQuery)
				}

				queryBody := httpQuery{
					Query:     newBody,
					Variables: variables,
				}

				sub := subQuery{
					Query:     lastQuery,
					Operation: string(operation.Operation),
					QueryName: field.Name,
					Body:      queryBody,
				}

				queries = append(queries, &sub)

			}

			if i > 0 {

				prevField := operation.SelectionSet[i-1].(*ast.Field)

				previosQuery := hq.Query[prevField.Position.Start:field.Position.Start]

				newBody := fmt.Sprintf(`%s %s { %s }`, operation.Operation, operation.Name, previosQuery)

				variables := make(map[string]interface{})

				variableStr := ""

				for i, queryVar := range prevField.Arguments {
					if val, ok := hq.Variables[queryVar.Name]; ok {
						variables[queryVar.Name] = val
					}

					if val, ok := operationMap[queryVar.Name]; ok {
						variableStr += "$" + queryVar.Name + ": " + val.Name

						if val.IsNullable {
							variableStr += "!"
						}

						if i != len(prevField.Arguments)-1 {
							variableStr += ","
						}
					}
				}

				if variableStr != "" {
					newBody = fmt.Sprintf(`%s %s(%s) { %s }`, operation.Operation, operation.Name, variableStr, previosQuery)

				}

				queryBody := httpQuery{
					Query:     newBody,
					Variables: variables,
				}

				sub := subQuery{
					Query:     previosQuery,
					Operation: string(operation.Operation),
					QueryName: prevField.Name,
					Body:      queryBody,
				}

				queries = append(queries, &sub)

			}

		}

	}

	return queries, nil

}

func (v *vertex) getMasterIntrospection() ([]byte, error) {
	byteData, err := v.schema.ToJSON()

	if err != nil {
		return nil, err
	}

	return byteData, nil
}
