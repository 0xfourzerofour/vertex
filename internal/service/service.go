package service

import (
	"embed"
	"errors"
	"govertex/internal/clients"
	"govertex/internal/graphql"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/cornelk/hashmap"
	"github.com/valyala/fasthttp"
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

var Client *fasthttp.Client

func LoadServices() error {
	readTimeout, _ := time.ParseDuration("5s")
	writeTimeout, _ := time.ParseDuration("s")
	maxIdleConnDuration, _ := time.ParseDuration("1h")

	Client = &fasthttp.Client{
		ReadTimeout:                   readTimeout,
		WriteTimeout:                  writeTimeout,
		MaxIdleConnDuration:           maxIdleConnDuration,
		NoDefaultUserAgentHeader:      true, // Don't send: User-Agent: fasthttp
		DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this
		DisablePathNormalizing:        true,
		// increase DNS cache time to an hour instead of default minute
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
	}

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

		for _, queryType := range svcIntrospection.Schema.QueryType.Fields {

			if field, ok := ServiceMap.GetStringKey(queryType.Name); ok {
				return errors.New(queryType.Name + "is already used and being send to " + field.(string))
			}

			ProxyMap.Insert(queryType.Name, svc.Url)

			if svc.Path != nil {
				ServiceMap.Insert(queryType.Name, *svc.Path)
			}
		}

		for _, mutationType := range svcIntrospection.Schema.MutationType.Fields {

			if field, ok := ServiceMap.GetStringKey(mutationType.Name); ok {
				return errors.New(mutationType.Name + "is already used and being send to " + field.(string))
			}

			ProxyMap.Insert(mutationType.Name, svc.Url)

			if svc.Path != nil {
				ServiceMap.Insert(mutationType.Name, *svc.Path)
			}
		}

		for _, subsriptionType := range svcIntrospection.Schema.SubscriptionType.Fields {

			if field, ok := ServiceMap.GetStringKey(subsriptionType.Name); ok {
				return errors.New(subsriptionType.Name + "is already used and being send to " + field.(string))
			}

			ProxyMap.Insert(subsriptionType.Name, svc.WS)

		}

	}

	return nil

}

func loadServicesFromDynamo() error {
	dynamoInput := dynamodb.QueryInput{
		TableName:              aws.String(os.Getenv("Vertex-Table")),
		KeyConditionExpression: aws.String("PK = :PK"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":PK": {
				S: aws.String("SERVICES"),
			},
		},
	}

	serviceList := []*ServiceItem{}

	dynamoOutput, err := clients.DynamoDB().Query(&dynamoInput)

	if err != nil {
		return err
	}

	if len(dynamoOutput.Items) != 0 {

		for _, item := range dynamoOutput.Items {

			service := ServiceItem{}

			err := dynamodbattribute.UnmarshalMap(item, &service)

			if err != nil {
				return err
			}

			serviceList = append(serviceList, &service)

		}

	}

	return nil
}
