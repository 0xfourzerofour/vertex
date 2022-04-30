### Vertex GraphQL service proxy ### 

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

## Considerations

 - There cannot be overlapping types throughout the graphql schemas
 - Services must have introspection turned on at the API level (Will work on a way to get around this)
 - This project is still early days so do not use in production

## TODO

 - Dockerize build for easy deployment and scalability
 - Logging
 - Load services from CLI
 - React UI portal to add services with API keys etc
