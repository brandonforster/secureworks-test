run:
	sqlite3 resolver.db '.read internal/sqlite/migrations/00_init.sql' &
	go run server.go

test:
	gofmt -l .
	[ "`gofmt -l .`" = "" ]
	go test ./...

docker_build:
	docker build --tag resolver:1.0 .

docker_run:
	docker run -p 8080:8080 --name resolver resolver:1.0

docker_clean:
	docker stop resolver
	docker rm resolver

helm:
	helm upgrade --install release helm-chart/ --values helm-chart/values.yaml
