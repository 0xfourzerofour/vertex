
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


