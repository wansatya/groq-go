.PHONY: build test lint

build:
	go build ./...

test:
	go test -v ./...

lint:
	golangci-lint run