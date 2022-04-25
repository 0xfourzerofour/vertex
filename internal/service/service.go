package service

import (
	"embed"
	"errors"
	"govertex/internal/graphql"
	"log"

	"github.com/cornelk/hashmap"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Services []*ServiceItem `json:"services"`
}
type ServiceItem struct {
	Url string `json:"url"`
	WS  string `json:"ws"`
}

var ServiceMap = hashmap.HashMap{}

//go:embed service-config.yml
var serviceEmbed embed.FS

func LoadServices() error {

	cfg := Config{}

	config, err := serviceEmbed.ReadFile("service-config.yml")

	if err != nil {
		log.Print(err)
	}

	err = yaml.Unmarshal(config, &cfg)

	if err != nil {
		log.Print(err)
	}

	for _, svc := range cfg.Services {

		svcIntrospection, err := graphql.GetIntrospectionSchema(svc.Url, "TESTTOKEN")

		if err != nil {
			return errors.New("Could not get introspection schema for " + svc.Url)
		}

		for _, queryType := range svcIntrospection.Schema.QueryType.Fields {
			ServiceMap.Insert(queryType.Name, svc.Url)
		}

	}

	return nil

}
