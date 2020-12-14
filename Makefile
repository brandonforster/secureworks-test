run:
	sqlite3 resolver.db
	go run server.go

test:
	gofmt -l .
	[ "`gofmt -l .`" = "" ]
	go test ./...

docker_build:
	docker build --tag resolver:1.0 .

docker_run:
	docker run -p 5000:5000 --name resolver resolver:1.0

docker_clean:
	docker stop resolver
	docker rm resolver
