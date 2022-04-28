package service

import (
	"embed"
	"errors"
	"govertex/internal/graphql"
	"log"

	"github.com/cornelk/hashmap"
	proxy "github.com/yeqown/fasthttp-reverse-proxy/v2"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Services []*ServiceItem `json:"services"`
}
type ServiceItem struct {
	Url  string  `json:"url"`
	WS   *string `json:"ws"`
	Path *string `json:"path"`
}

var ServiceMap = hashmap.HashMap{}
var ProxyMap = hashmap.HashMap{}

//go:embed service-config.yml
var serviceEmbed embed.FS

func LoadServices() error {

	cfg := Config{}

	config, err := serviceEmbed.ReadFile("service-config.yml")

	if err != nil {
		return err
	}

	err = yaml.Unmarshal(config, &cfg)

	if err != nil {
		log.Print(err)
	}

	for _, svc := range cfg.Services {

		serviceUri := svc.Url
		if svc.Path != nil {
			serviceUri += *svc.Path
		}

		svcIntrospection, err := graphql.GetIntrospectionSchema("https://"+serviceUri, "TESTTOKEN")

		if err != nil {
			return errors.New("Could not get introspection schema for " + svc.Url)
		}

		urlProxy := proxy.NewReverseProxy(svc.Url)

		for _, queryType := range svcIntrospection.Schema.QueryType.Fields {

			if field, ok := ServiceMap.GetStringKey(queryType.Name); ok {
				return errors.New(queryType.Name + "is already used and being send to " + field.(string))
			}

			ProxyMap.Insert(queryType.Name, urlProxy)

			if svc.Path != nil {
				ServiceMap.Insert(queryType.Name, *svc.Path)
			}
		}

	}

	return nil

}
