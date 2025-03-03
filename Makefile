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
	go build github.com/darkturo/slocc/cmd/slocc
.PHONY: build

test:
	go test ./...
.PHONY: test
