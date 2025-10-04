# KubeVirt MCP Server Agent Instructions

This document provides instructions for AI agents working on the KubeVirt MCP server project.

## Project Overview

This is a simple MCP (Model Context Protocol) server implementation for the KubeVirt project, a virtualization extension for Kubernetes. It uses the official Go SDK for MCP.

### Project Structure
- `main.go`: MCP server setup and registration
- `pkg/client/`: Shared KubeVirt client utilities
- `pkg/tools/`: MCP tool handlers for VM operations
- `pkg/resources/`: MCP resource handlers for structured data access
- `scripts/kubevirtci.sh`: Script for managing local kubevirtci development environment
- `scripts/sync.sh`: Script for building and running MCP server locally with kubevirtci access
- `Makefile`: Build automation and development tasks
- `go.mod` / `go.sum`: Go module dependencies

## Development Setup

For functional testing with a real KubeVirt cluster, use the kubevirtci integration.

- To start a local kubevirtci cluster with KubeVirt:
  ```bash
  make cluster-up
  ```
- To stop the cluster:
  ```bash
  make cluster-down
  ```

The `scripts/kubevirtci.sh` script handles downloading and setting up kubevirtci, starting a local Kubernetes cluster with KubeVirt and CDI, and configuring the environment for testing.

## Building and Running

### Building the server
Use the Makefile to build the binary:
```bash
make build
```
Alternatively, use the Go command:
```bash
go build -o kubevirt-mcp-server .
```

### Running the server
To build and run the server locally with access to the kubevirtci cluster:
```bash
make cluster-sync
```
The `scripts/sync.sh` script handles building the server, starting it with the correct `KUBECONFIG`, and logging output to `/tmp/kubevirt-mcp-server.log`.

### Other Make Targets
- `make clean`: Clean build artifacts
- `make deps`: Download and tidy dependencies
- `make run`: Build and run the server
- `make help`: Show all available make targets

## Testing and Code Quality

### Running Tests
The project uses the Ginkgo testing framework.
- To run unit tests:
  ```bash
  make test
  ```
- To generate a test coverage report:
  ```bash
  make coverage
  ```
- To run functional tests against the MCP server (requires `make cluster-up`):
  ```bash
  make test-functional
  ```

### Code Quality Checks
The project uses `golangci-lint` for linting and standard Go tools for formatting and vetting.
- To format Go code:
  ```bash
  make fmt
  ```
- To run `go vet`:
  ```bash
  make vet
  ```
- To run the linter:
  ```bash
  make lint
  ```
- To run all checks (fmt, vet, lint, test):
  ```bash
  make check
  ```

### Test Structure
- `pkg/client/client_test.go`: Tests for KubeVirt client creation
- `pkg/tools/tools_test.go`: Tests for MCP tool handlers
- `pkg/resources/resources_test.go`: Tests for MCP resource handlers
- `tests/functional/`: Functional tests for MCP server stdio communication

## Available Tools and Resources

The server exposes the following capabilities via the Model Context Protocol (MCP).

### MCP Tools
- `list_vms`: List virtual machine names in a namespace.
- `start_vm`, `stop_vm`, `restart_vm`: Manage VM state.
- `pause_vm`, `unpause_vm`: Pause and unpause a VM.
- `create_vm`: Create a new VM. Supports OS name shortcuts for container disks.
- `delete_vm`: Delete a VM.
- `patch_vm`: Apply a JSON merge patch to a VM.
- `list_instancetypes`, `get_instancetype`: Manage instance types.
- `get_preference`: Get preference details.
- `get_vm_instancetype`, `get_vm_status`, `get_vm_conditions`, `get_vm_phase`, `get_vm_disks`: Get VM details.

### MCP Prompts
- `describe_vm`: Provide a comprehensive VM description.
- `troubleshoot_vm`: Diagnose potential VM issues.
- `health_check_vm`: Perform a quick VM health check.

### MCP Resources
- `kubevirt://{namespace}/vms`: List of VMs.
- `kubevirt://{namespace}/vm/{name}`: Complete VM specification.
- `kubevirt://{namespace}/vm/{name}/status`: VM status.
- `kubevirt://{namespace}/vmis`: List of VMIs.
- `kubevirt://{namespace}/vmi/{name}`: Complete VMI specification.
- And many more for data volumes, instance types, preferences, etc.

## Contribution Guidelines

### Commit Messages
This project follows the [Conventional Commits v1.0.0](https://www.conventionalcommits.org/en/v1.0.0/) specification.
- **Types**: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `revert`.
- **Example**: `feat(tools): add VM restart functionality`

### AI Assistance Attribution
When using AI assistance, add the `Assisted-By` trailer to your commit message.
```
feat(resources): add VM instance filtering

Add support for filtering VMs by status and labels in resource handlers.

Assisted-By: Claude <noreply@anthropic.com>
```

## Documentation

**CRITICAL REQUIREMENT**: When making ANY changes to the project (code, features, tools, etc.), you MUST update BOTH `AGENTS.md` and `README.md`. This ensures documentation remains current for both agents and human developers.