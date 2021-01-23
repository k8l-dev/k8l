generate:
	npx @openapitools/openapi-generator-cli generate -g go-gin-server -p apiPath=api -p packageName=api -i static/openapi.yaml -o go

build: 
	go build -o k8l go/main.go

dev:
	nodemon --exec go run -tags "fts5" go/main.go -listen :9099 -data ./data -verbose --signal SIGTERM

clean:
	rm -rf ./k8l