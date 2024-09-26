# Makefile for NVCF CLI
BINARY_NAME := nvcf

## Default target: build the project
all: check install

## Install the binary
install:
	go install

## Build the binary
build:
	go build -o $(BINARY_NAME)

## Clean up build artifacts
clean:
	go clean
	rm -f $(BINARY_NAME)

## Run tests
test:
	go test ./...

## Run all checks: fmt, vet, lint, and test
check: fmt vet lint test

## Format the code
fmt:
	gofmt -s -w .

## Run go vet
vet:
	go vet ./...

## Run linter
lint: deps-lint
	golangci-lint run

## Quick rebuild: clean and build
q: clean build

## Install linter dependency
deps-lint:
	@command -v golangci-lint > /dev/null || (go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0)

## Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  %-20s %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

## Build docs
docs:
	go run . docs

## Clean up docs dir
cleandocs:
	rm -rf ./docs

dry-release:
	goreleaser release --snapshot --clean

release:
	goreleaser release --clean

.PHONY: all build clean test check fmt vet lint q deps-lint help

