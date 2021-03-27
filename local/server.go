package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/entegral/aboutme/aws"
	"github.com/entegral/aboutme/graph/generated"
	"github.com/entegral/aboutme/graph/resolvers"
)

const defaultPort = "8989"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	
	LocalServer := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers.Clients}))
	aws.AddCustomClientConnection("us-west-2", "default")
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", LocalServer)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}