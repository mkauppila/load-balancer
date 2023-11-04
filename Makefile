test:
	go test ./...
.PHONY: test

lint: 
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2 run ./...
.PHONY: lint

run:
	go run cmd/loadbalancer/main.go
.PHONY: run

httpserver:
	HTTP_PORT=12000 go run cmd/httpserver/main.go
.PHONY: httpserver

docker-build-loadbalancer:
	docker build -f loadbalancer.docker -t loadbalancer .
.PHONY: docker-build-loadbalancer

docker-build-httpserver:
	docker build -f httpserver.docker -t httpserver --build-arg PORT=12000 .
.PHONY: docker-build-httpserver
