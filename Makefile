test:
		go test ./...
.PHONY: test

lint: 
		go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2 run ./...
.PHONY: lint

run:
		go run cmd/main.go
.PHONY: run

