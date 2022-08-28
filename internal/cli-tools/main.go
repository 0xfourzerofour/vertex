package cli_tools

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
)

type vertexData struct {
	ServiceName string            `json:"serviceName"`
	ServiceUrl  string            `json:"serviceUrl"`
	Schema      string            `json:"schema"`
	QueryMap    map[string]string `json:"queryMap"`
}

func loadSchema(path string) ([]byte, error) {
	file, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func saveVertexFile(vertexFile []byte, outputPath string) error {

	err := ioutil.WriteFile(outputPath, vertexFile, 0644)

	if err != nil {
		return err
	}

	return nil
}

func GenerateVertex(inputFile, outputFile, serviceName, serviceUrl string) error {
	schema, err := loadSchema(inputFile)

	if err != nil {
		return err
	}

	serviceMap, err := generateMap(string(schema), serviceUrl)

	if err != nil {
		return err
	}

	schemaData := vertexData{
		ServiceName: serviceName,
		ServiceUrl:  serviceUrl,
		Schema:      string(schema),
		QueryMap:    serviceMap,
	}

	schemaBytes, err := json.MarshalIndent(schemaData, "", "\t")

	if err != nil {
		return err
	}

	err = saveVertexFile(schemaBytes, outputFile)

	if err != nil {
		return err
	}

	return nil

}

func generateMap(schemaString, serviceUrl string) (map[string]string, error) {

	astSource := ast.Source{
		Input: schemaString,
	}

	schema, err := gqlparser.LoadSchema(&astSource)

	if err != nil {
		return nil, err
	}

	schemaMap := map[string]string{}

	if schema.Mutation != nil {
		for _, mutation := range schema.Mutation.Fields {
			schemaMap[mutation.Name] = serviceUrl
		}
	}

	if schema.Query != nil {
		for _, query := range schema.Query.Fields {
			schemaMap[query.Name] = serviceUrl
		}
	}

	if schema.Subscription != nil {
		for _, sub := range schema.Subscription.Fields {
			schemaMap[sub.Name] = serviceUrl
		}
	}

	return schemaMap, nil

}
