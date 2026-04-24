BINARY_NAME := linear-cli
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS := -s -w \
	-X github.com/enolalab/linear-cli/cmd.version=$(VERSION) \
	-X github.com/enolalab/linear-cli/cmd.commit=$(COMMIT) \
	-X github.com/enolalab/linear-cli/cmd.date=$(DATE) \
	-X github.com/enolalab/linear-cli/internal/api.Version=$(VERSION)

.PHONY: build clean test lint install

## build: Build the binary
build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) .

## install: Install to $GOPATH/bin
install:
	go install -ldflags "$(LDFLAGS)" .

## test: Run tests
test:
	go test -v -race ./...

## lint: Run linter
lint:
	golangci-lint run ./...

## clean: Remove build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -rf dist/

## help: Show this help
help:
	@echo "Available targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  /'
