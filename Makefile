.PHONY: build test install clean lint bench help

VERSION := v0.1.0
BINARY_NAME := htmlx
OUTPUT_DIR := bin

help:
	@echo "htmlx - HTML to Markdown Converter"
	@echo ""
	@echo "Available targets:"
	@echo "  build       Build the htmlx binary"
	@echo "  test        Run all tests"
	@echo "  install     Install htmlx to \$$GOPATH/bin"
	@echo "  clean       Remove build artifacts"
	@echo "  lint        Run golangci-lint"
	@echo "  bench       Run benchmarks"
	@echo "  help        Show this help message"

build:
	@mkdir -p $(OUTPUT_DIR)
	@echo "Building $(BINARY_NAME)..."
	@go build -v -o $(OUTPUT_DIR)/$(BINARY_NAME) ./cmd/htmlx

test:
	@echo "Running tests..."
	@go test -v -race -cover ./...

install: build
	@echo "Installing $(BINARY_NAME)..."
	@go install ./cmd/htmlx

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(OUTPUT_DIR)/
	@go clean

lint:
	@echo "Running linter..."
	@golangci-lint run ./...

bench:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

.PHONY: version
version:
	@echo "htmlx $(VERSION)"
