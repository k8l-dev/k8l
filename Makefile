generate:
	npx @openapitools/openapi-generator-cli generate -g go-gin-server -p apiPath=api -p packageName=api -i static/openapi.yaml -o go


build:
	go build -o bin/main main.go

dev:
	go run -tags "fts5" go/main.go -listen :9099 -data ./data

clean:
	rm -rf ./k8l