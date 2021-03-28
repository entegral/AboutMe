.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/lambda lambda/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy

generate:
	go run github.com/99designs/gqlgen generate 

local: clean build
	go run local/server.go

deploy: clean build
	sls deploy --verbose
