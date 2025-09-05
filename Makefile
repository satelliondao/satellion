# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=bin/sat
BINARY_UNIX=$(BINARY_NAME)_unix

# Build the application
.PHONY: build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v .

# Build for multiple platforms
.PHONY: build-all
build-all:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/sat_linux_amd64 -v .
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o bin/sat_windows_amd64.exe -v .
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o bin/sat_darwin_amd64 -v .
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o bin/sat_darwin_arm64 -v .

# Clean build artifacts
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f bin/sat_*

# Run tests
.PHONY: test
test:
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

# Download dependencies
.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run the application
.PHONY: run
run:
	$(GOCMD) run main.go

# Format code
.PHONY: fmt
fmt:
	$(GOCMD) fmt ./...

# Lint code
.PHONY: lint
lint:
	golangci-lint run

# Install development tools
.PHONY: install-tools
install-tools:
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint

# Development build with race detection
.PHONY: build-dev
build-dev:
	$(GOBUILD) -race -o ./bin/$(BINARY_NAME) -v .

# Install the binary to GOPATH/bin
.PHONY: install
install:
	$(GOCMD) install -v ./...

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  build-all    - Build for multiple platforms"
	@echo "  build-dev    - Build with race detection"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  deps         - Download and tidy dependencies"
	@echo "  run          - Run the application"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code (requires golangci-lint)"
	@echo "  install-tools- Install development tools"
	@echo "  install      - Install binary to GOPATH/bin"
	@echo "  help         - Show this help message"

# Default target
.DEFAULT_GOAL := build
