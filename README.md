# resolver    
 Please see the [specs sheet for a better understanding of what this program attempts to do.](requirements.md)    
    
### Building the Application    
 This section assumes that you have the Go toolchain and SQLite3 installed previously and on your system's path.    
    
Simply run `make run` and a script will start that will launch the database engine, install dependencies, and launch the server. From there, you can open up your browser and navigate to `localhost:8080/` where you'll find the GraphQL playground to get started interacting with the system.    
      
### Using the Application    
 The two primary ways of interacting with the application are the `enqueue` function and the `getIPDetails` function. `enqueue` should be thought of as a "permanent" way to look up data, where we run a query for an IP address against the blocklist and store that result in the database. `getIPDetails` should be thought of as a "transient" way to look up data, as once the query completes the object that was returned from the function is not stored anywhere.    
       
For help using GraphQL or making queries with it, please consult your favorite search engine or tutorial provider.    
       
### Alternative Deployments    
 In addition to "normal" local build and development, this application supports Docker deployments. Simply run `make docker_build`, `make docker_run` to get started with Docker. You'll need to run `make docker_clean` between builds of the application. You can customize the port the app runs on by editing the `PORT` environment variable in the Dockerfile, but remember that you'll need to update either the `Makefile` or your run command to use the new port you specified.    
  
To run the app in Kubernetes:  
 - Use the Helm chart to deploy it. Run `make helm` to kick off the deployment.  
 - Next, you need to configure your network to be able to interact with the application running inside the pod. Run `kubectl get pods` and take node of the name of the pod. It should be some variation of `resolver-chart-` with a string of alphanumeric digits behind the dash.
 - Run the following command: `kubectl port-forward {POD NAME} 8080:8080`, replacing `{POD NAME}` with the name you looked up in the previous step. This will expose port 8080 in the Kubernetes cluster to your local machine.
 - You can now use your browser and access the application at `localhost:8080`.
      
If your system doesn't have access to the program `make` you can run all the same commands yourself. Refer to the [Makefile](Makefile) and run each line of the script you'd like to use one at a time.    
    
## Dependency Justification    
 `github.com/go-chi/chi`    
 This router was chosen because it seems to integrate the best with the gqlgen framework. I had little trouble adding my implementation of basic authentication, and implementing it in an idiomatic-to-gqlgen way. I liked how implementing my middleware was a one and done affair.    
    
`github.com/google/uuid`    
 UUID is one of the defined fields on the spec, and this is the most widely used UUID library and one I've used in the past. It is elegant and concise and never feels lacking in any methods that you might want.    
     
``` github.com/jmoiron/sqlx github.com/mattn/go-sqlite3 ```   
These two libraries are what I've used in the past for handling SQLite in Go. sqlx, while an abstraction, makes operating with SQL of any flavor very easy and configuring it to use SQLite is just a matter of importing `go-sqlite3`.
