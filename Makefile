.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/lambda lambda/main.go

clean:
	rm -rf ./bin

local: clean build
	go run local/server.go

deploy: clean build
	sls deploy --verbose
