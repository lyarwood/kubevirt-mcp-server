# KubeVirt MCP Server Makefile

BINARY_NAME := kubevirt-mcp-server
BUILD_DIR := .
GO_FILES := $(shell find . -type f -name '*.go')

.PHONY: build clean test fmt vet deps help

# Default target
all: build

# Build the binary
build: $(BINARY_NAME)

$(BINARY_NAME): $(GO_FILES)
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) .
	@echo "Build complete!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	go clean
	@echo "Clean complete!"

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Format Go code
fmt:
	@echo "Formatting Go code..."
	go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Run the server (for development)
run: build
	@echo "Starting $(BINARY_NAME)..."
	./$(BINARY_NAME)

# Check code quality
check: fmt vet test
	@echo "Code quality checks complete!"

# Show help
help:
	@echo "Available targets:"
	@echo "  build    - Build the binary (default)"
	@echo "  clean    - Clean build artifacts"
	@echo "  test     - Run tests"
	@echo "  fmt      - Format Go code"
	@echo "  vet      - Run go vet"
	@echo "  deps     - Download and tidy dependencies"
	@echo "  run      - Build and run the server"
	@echo "  check    - Run fmt, vet, and test"
	@echo "  help     - Show this help message"