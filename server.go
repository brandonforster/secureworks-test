package main

import (
	"context"
	"crypto/subtle"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"

	"github.com/brandonforster/resolver/graph"
	"github.com/brandonforster/resolver/graph/generated"
)

const PORT = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = PORT
	}

	router := chi.NewRouter()
	router.Use(basicAuth())

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/graphql", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// literally basic auth
func basicAuth() func(http.Handler) http.Handler {
	username := "secureworks"
	password := "supersecret"

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()
			var ctx context.Context

			if !ok ||
				subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 ||
				subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
				ctx = context.WithValue(r.Context(), "isAuth", false)
			} else {
				ctx = context.WithValue(r.Context(), "isAuth", true)
			}

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
