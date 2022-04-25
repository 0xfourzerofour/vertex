### Vertex GraphQL service proxy ### 

This Project is inspired by current pitfalls that I have come across at work with monolithic graphql schemas.
Vertex aims to solve this issue by allowing a single graphql endpoint to many downstream services by parsing 
the query body and matching the query to a service. 

# Test Query 1 #

```

// countries.trevorblades.com

{
  languages {
        code
      }
    }
  }
}


```


# Test Query 2 #

```

// www.universe.com/graphql

{
  episodes {
    info {
      next
    }
  }
}

```

