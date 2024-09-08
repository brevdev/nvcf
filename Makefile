# Makefile for NVCF CLI

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
BINARY_NAME=nvcf
GOFILES=$(shell find . -name '*.go' -not -path "./vendor/*")

# Linting
GOLINT=golangci-lint

.PHONY: all build clean test vet lint fmt q

all: build

build:
	$(GOBUILD) -o $(BINARY_NAME)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

test:
	$(GOTEST) ./...

vet:
	$(GOVET) ./...

lint:
	$(GOLINT) run

fmt:
	gofmt -s -w $(GOFILES)

# Installs golangci-lint
install-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2

# Runs all checks
check: fmt vet lint test

# Builds for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-linux-amd64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-darwin-amd64
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-windows-amd64.exe

q:
	make clean
	make build