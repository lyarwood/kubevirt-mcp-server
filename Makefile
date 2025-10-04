# KubeVirt MCP Server Makefile

BINARY_NAME := kubevirt-mcp-server
BUILD_DIR := .
GO_FILES := $(shell find . -type f -name '*.go')
IMAGE_NAME := kubevirt-mcp-server
IMAGE_TAG := latest

.PHONY: build clean test fmt vet deps help image cluster-up cluster-down cluster-sync test-functional

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

# Run unit tests only
test:
	@echo "Running unit tests with Ginkgo..."
	ginkgo --randomize-all --randomize-suites --fail-on-pending --cover --trace -v pkg/...

# Generate test coverage report for unit tests
coverage:
	@echo "Generating test coverage report for unit tests..."
	ginkgo --randomize-all --randomize-suites --fail-on-pending --cover --trace --coverprofile=coverage.out pkg/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Format Go code
fmt:
	@echo "Formatting Go code..."
	go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Run linter
lint:
	@echo "Running golangci-lint..."
	golangci-lint run --timeout=5m

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Run the server (for development)
run: build
	@echo "Starting $(BINARY_NAME)..."
	./$(BINARY_NAME)

# Run all quality checks
check: fmt vet lint test
	@echo "Code quality checks complete!"

# Build container image
image:
	@echo "Building container image $(IMAGE_NAME):$(IMAGE_TAG)..."
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .
	@echo "Container image built successfully!"

# Start kubevirtci cluster
cluster-up:
	@echo "Starting kubevirtci cluster..."
	./scripts/kubevirtci.sh up
	@echo "Cluster started! Use 'make cluster-down' to stop."

# Stop kubevirtci cluster
cluster-down:
	@echo "Stopping kubevirtci cluster..."
	./scripts/kubevirtci.sh down
	@echo "Cluster stopped!"

# Build and run MCP server locally with kubevirtci access
cluster-sync:
	@echo "Building and starting MCP server locally..."
	./scripts/sync.sh
	@echo "MCP server started successfully!"

# Run functional tests against MCP server
test-functional: build
	@echo "Running functional tests..."
	ginkgo run --randomize-all --randomize-suites --trace -v tests/functional
	@echo "Functional tests complete!"

# Show help
help:
	@echo "Available targets:"
	@echo "  build     - Build the binary (default)"
	@echo "  clean     - Clean build artifacts"
	@echo "  test      - Run tests with Ginkgo framework"
	@echo "  coverage  - Generate test coverage report"
	@echo "  fmt       - Format Go code"
	@echo "  vet       - Run go vet"
	@echo "  lint      - Run golangci-lint"
	@echo "  deps      - Download and tidy dependencies"
	@echo "  run       - Build and run the server"
	@echo "  check     - Run fmt, vet, lint, and test"
	@echo "  image     - Build container image"
	@echo "  cluster-up   - Start kubevirtci cluster for testing"
	@echo "  cluster-down - Stop kubevirtci cluster"
	@echo "  cluster-sync - Build and run MCP server locally with kubevirtci access"
	@echo "  test-functional - Run functional tests against MCP server"
	@echo "  help      - Show this help message"