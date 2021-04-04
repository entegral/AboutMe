package main

import (
	"context"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/99designs/gqlgen/handler"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/handlerfunc"
	"github.com/entegral/aboutme/errors"
	"github.com/entegral/aboutme/graph/generated"
	"github.com/entegral/aboutme/graph/resolvers"
	"github.com/sirupsen/logrus"
)
	


var ApiGatewayAdapter *handlerfunc.HandlerFuncAdapter
var ApiGatewayPlayground *handlerfunc.HandlerFuncAdapter

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

func init() {
	schema := generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers.Clients})
	gqlHandler := handler.GraphQL(schema)
	ApiGatewayAdapter = handlerfunc.New(gqlHandler)
	ApiGatewayPlayground = handlerfunc.New(playground.Handler("Playground", "/"))
}

// Handler receives the event from lambda.Start and processes the request
func Handler(ctx context.Context, req Request) (Response, error) {
	logrus.WithFields(logrus.Fields{
		"body": req.Body,
		"path": req.Path,
		"method": req.HTTPMethod,
	}).Debugf("request: %+v", req)

	if req.Path =="/" && req.HTTPMethod == "GET" {
		logrus.Debugln("inside get route, querying playground")
		rsp, err := ApiGatewayPlayground.ProxyWithContext(ctx, events.APIGatewayProxyRequest(req))
		if errors.Warn("main.Handler.GET", err) {
			return Response{
				StatusCode: 500,
				Body: "internal service error:" + err.Error(),
			}, err
		}
		return Response(rsp), nil
	} else if req.Path == "/query" || req.HTTPMethod == "POST" {

		rsp, err := ApiGatewayAdapter.ProxyWithContext(ctx,events.APIGatewayProxyRequest(req))
		if errors.Warn("main.Handler.POST", err) {
			return Response{
				StatusCode: 500,
				Body: string(rsp.Body),
			}, nil 
		}
		return Response(rsp), nil
	}
	
	return Response{
		StatusCode: 400,
		Body: "endpoint only supports GET request on '/' path and POST request on '/query' path",
	}, nil
}

func main() {
	lambda.Start(Handler)
}
