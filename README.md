# resolver

Please see the [specs sheet for a better understanding of what this program attempts to do.](requirements.md)

### Building the Application

This section assumes that you have the Go toolchain and SQLite3 installed previously and on your system's path.

Simply run `make run` and a script will start that will launch the database engine, install dependencies, and launch
 the server. From there, you can open up your browser and navigate to `localhost:8080/` where you'll find the GraphQL
  playground to get started interacting with the system.
  
### Using the Application

The two primary ways of interacting with the application are the `enqueue` function and the `getIPDetails` function.
 `enqueue` should be thought of as a "permanent" way to look up data, where we run a query for an IP address against
  the blocklist and store that result in the database. `getIPDetails` should be thought of as a "transient" way to
   look up data, as once the query completes the object that was returned from the function is not stored anywhere.
   
   For help using GraphQL or making queries with it, please consult your favorite search engine or tutorial provider.
   
### Alternative Deployments

In addition to "normal" local build and development, this application supports Docker and Kubernetes deployments.
 Simply run `make docker_build`, `make docker_run` to get started with Docker and `make k8s` to deploy the
  application to a running Kubernetes cluster that you have previously set up. 

## Dependency Justification
`github.com/99designs/gqlgen`

This is used because the spec sheet recommends it.
	
`github.com/go-chi/chi`

This is used because the gqlgen tutorial used it to great effect, and deviating from the recommendation of a
 recommendation seemed unwise.

`github.com/google/uuid`

UUID is one of the defined fields on the spec, and this is the most widely used UUID library and one I've used in the
 past.
 

    github.com/jmoiron/sqlx
    github.com/mattn/go-sqlite3
    
These two libraries are what I've used in the past for handling SQLite in Go. sqlx, while an abstraction, makes
 operating with SQL of any flavor very easy and configuring it to use SQLite is just a matter of importing go-sqlite3.

    github.com/vektah/gqlparser/v2
    
This is imported automatically by gqlgen.
