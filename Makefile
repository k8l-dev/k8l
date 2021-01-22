api/routers.go:
	npx @openapitools/openapi-generator-cli generate -g go-gin-server -p apiPath=api -p packageName=api -i static/openapi.yaml -o go


build: api/routers.go
	go build -o bin/main main.go

dev:
	nodemon --exec go run -tags "fts5" go/main.go -listen :9099 -data ./data --signal SIGTERM

clean:
	rm -rf ./k8l