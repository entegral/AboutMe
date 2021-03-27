package main

import (
	"context"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/99designs/gqlgen/handler"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/entegral/aboutme/errors"
	"github.com/entegral/aboutme/graph/generated"
	"github.com/entegral/aboutme/graph/resolvers"
)


var ApiGatewayAdapter lambda.Handler
var ApiGatewayPlayground lambda.Handler

// Response is of type APIGatewayProxyResponse since we're leveraging the
type Response events.APIGatewayProxyResponse

func init() {
	schema := generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers.Clients})
	gqlHandler := handler.GraphQL(schema)
	// ApiGatewayAdapter = handlerfunc.New(gqlHandler)
	ApiGatewayAdapter = lambda.NewHandler(gqlHandler)
	ApiGatewayPlayground = lambda.NewHandler(playground.Handler("Playground", "/"))
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (Response, error) {

	if req.Path =="/" && req.HTTPMethod == "GET" {
		rsp, err := ApiGatewayPlayground.Invoke(ctx, []byte(req.Body))
		if errors.Warn("main.Handler.GET", err) {
			return Response{
				StatusCode: 500,
				Body: "internal service error",
			}, err
		}
		return Response{
			StatusCode: 500,
			Body: string(rsp),
		}, nil
	} else if req.Path == "/query" || req.HTTPMethod == "POST" {

		rsp, err := ApiGatewayAdapter.Invoke(ctx, []byte(req.Body))
		if errors.Warn("main.Handler.POST", err) {
			return Response{
				StatusCode: 500,
				Body: string(rsp),
			}, nil 
		}
		return Response{
			StatusCode: 200,
			Body: string(rsp),
		}, nil
	}
	
	return Response{
		StatusCode: 500,
		Body: "endpoint only supports GET request on '/' path and POST request on '/query' path",
	}, nil
}

func main() {
	lambda.Start(Handler)
}
