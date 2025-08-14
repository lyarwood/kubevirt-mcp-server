# kubevirt-mcp-server

A simple Model Context Protocol server for KubeVirt.

## Architecture

The project is organized into modular packages:

- `main.go` - MCP server setup and registration
- `pkg/client/` - Shared KubeVirt client utilities
- `pkg/tools/` - MCP tool handlers for VM operations
- `pkg/resources/` - MCP resource handlers for structured data access
- `Makefile` - Build automation and development tasks

## Features

### MCP Tools
- `list_vms` - List virtual machine names in a namespace
- `start_vm` - Start a virtual machine
- `stop_vm` - Stop a virtual machine
- `list_instancetypes` - List available instance types
- `get_vm_instancetype` - Get instance type for a VM

### MCP Resources
- `kubevirt://{namespace}/vms` - JSON list of VMs with summary info
- `kubevirt://{namespace}/vm/{name}` - Complete VM specification
- `kubevirt://{namespace}/vmis` - JSON list of VMIs with runtime info
- `kubevirt://{namespace}/vmi/{name}` - Complete VMI specification

## Building

```bash
# Using Makefile (recommended)
make build

# Or directly with go
go build -o kubevirt-mcp-server .
```

## Development

### Available Make Targets
- `make build` - Build the binary (default)
- `make clean` - Clean build artifacts  
- `make test` - Run tests with Ginkgo framework
- `make coverage` - Generate test coverage report
- `make fmt` - Format Go code
- `make vet` - Run go vet
- `make deps` - Download and tidy dependencies
- `make run` - Build and run the server
- `make check` - Run fmt, vet, and test
- `make help` - Show help message

### Testing

The project uses the [Ginkgo](https://github.com/onsi/ginkgo) testing framework with [Gomega](https://github.com/onsi/gomega) assertions:

```bash
# Run all tests
make test

# Generate test coverage report
make coverage
```

## Demo

This short demo uses mcp-cli as a bridge between the kubevirt-mcp-server and LLM.

The model used by the demo is llama3.2 running locally under ollama.

![demo](demo.gif)

## Links 

- https://www.anthropic.com/news/model-context-protocol
- https://github.com/mark3labs/mcp-go
- https://github.com/chrishayuk/mcp-cli
