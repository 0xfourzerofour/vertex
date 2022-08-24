### Vertex GraphQL Api Gateway ### 

This Project is inspired by current pitfalls that I have come across at work with monolithic graphql schemas.
Vertex aims to solve this issue by allowing a single graphql endpoint to many downstream services by parsing 
the query body and matching the query to a service. 

[![](https://mermaid.ink/img/pako:eNp9j00KwjAQRq8SZt1eIAtB2qgLF9UGFZouQhNtsElKTMHS9u6m_iwEcRbD8L3HMDNAZYUEDBfH2xrRlBkUalkkjZLGlyiOF-OG0iyaWx4d8xElw4HsKTlNLzd5Opmz935EabGeF-225Q9I_sHVF4QItHSaKxFuG-aEga-llgxwGAV3VwbMTMHrWsG9JEJ56wCfeXOTEfDO27w3FWDvOvmRUsXDn_ptTQ8RuVGX)](https://mermaid-js.github.io/mermaid-live-editor/edit#pako:eNp9j00KwjAQRq8SZt1eIAtB2qgLF9UGFZouQhNtsElKTMHS9u6m_iwEcRbD8L3HMDNAZYUEDBfH2xrRlBkUalkkjZLGlyiOF-OG0iyaWx4d8xElw4HsKTlNLzd5Opmz935EabGeF-225Q9I_sHVF4QItHSaKxFuG-aEga-llgxwGAV3VwbMTMHrWsG9JEJ56wCfeXOTEfDO27w3FWDvOvmRUsXDn_ptTQ8RuVGX)

## Service Config

services are added via the `internal/service/service-config.yml`

```

services:
    - 
        url: "yourgraphqlservice.com"
        ws: "ws://yourgraphqlwebsocket"
        path: "/graphql" //Optional

```

## Example Query

```

query vertex($id: ID!) {

    //fruists api

    fruit(id: $id) {
        description     
    }
    
    //countries api    
    
    languages {
        code
    }
    
    //rick and morty api
    
    characters(page: 2, filter: { name: "Morty" }) {
      info {
        count
      }
      results {
        name
      }
    }
}

```

## Exciting findings

 - requests can be parsed within 0.01-0.1ms which will only get faster after I implement request caching

## Considerations

 - There cannot be overlapping types throughout the graphql schemas (Might be an option for query polymorphism based on different variables)
 - Services must have introspection turned on at the API level (I am working around this using a ci tool to publish schemas)
 - This project is still early days so do not use in production

## TODO

 - Create Nextjs platform for hosted graphql schemas
 - Create cli tool to publish schemas from service CI
 - Add terraform deployment to autoscaling group behind a load balancer
 - create cli tool to deploy schemas from service ci to remove need for introspection
 - Load all graphs into one schema for playground introspection
 - Add Web socket functionality for subsriptions
 - cache full introspection graph in memory for fast retrieval
 - Cache hash of request body for faster proxying
 - Load and save services from DynamoDB
 - Allow passing variables to query/mutations
 - Dockerize build for easy deployment and scalability
 - Logging
 - Load services from CLI
