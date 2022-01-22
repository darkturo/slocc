.DEFAULT_GOAL := build

fmt:
	go fmt ./...
.PHONY: fmt

lint: fmt
	golangci-lint run
.PHONY: lint

vet: lint
	go vet ./...
.PHONY: vet

build: vet
	go build -o slocc main.go 
.PHONY: build

test:
	go test ./...
.PHONY: test
