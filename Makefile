.PHONY: build install clean test run help

BINARY_NAME=nco
BUILD_DIR=build
VERSION=$(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.1.0")
REVISION=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILT_AT=$(shell date -u '+%Y-%m-%dT%H:%M:%S')

BUILD_VARS=github.com/ncobase/cli/version
LDFLAGS=-s -w -X $(BUILD_VARS).Version=$(VERSION) -X $(BUILD_VARS).Revision=$(REVISION) -X $(BUILD_VARS).BuiltAt=$(BUILT_AT)
GOFLAGS=-trimpath -ldflags="$(LDFLAGS)"

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the CLI binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

install: ## Install the CLI binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	go install $(GOFLAGS) .

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	go clean

test: ## Run tests
	@echo "Running tests..."
	go test ./...

run: ## Run the CLI locally
	@echo "Running $(BINARY_NAME)..."
	go run . $(ARGS)

build-all: ## Build for multiple platforms
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

.DEFAULT_GOAL := help
