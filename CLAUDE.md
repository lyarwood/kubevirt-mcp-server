# KubeVirt MCP server

This is a simple MCP server implementation for the KubeVirt project, a virtualization extension for Kubernetes.

## Project Structure

- `main.go` - MCP server setup and registration
- `pkg/client/` - Shared KubeVirt client utilities
- `pkg/tools/` - MCP tool handlers for VM operations
- `pkg/resources/` - MCP resource handlers for structured data access
- `scripts/kubevirtci.sh` - Script for managing local kubevirtci development environment
- `scripts/sync.sh` - Script for building and running MCP server locally with kubevirtci access
- `Makefile` - Build automation and development tasks
- `server_config.json` - MCP server configuration for client integration
- `go.mod` / `go.sum` - Go module dependencies

## Features

### MCP Tools
- `list_vms` - List virtual machine names in a namespace
- `start_vm` - Start a virtual machine
- `stop_vm` - Stop a virtual machine
- `restart_vm` - Restart a virtual machine (handles both running and stopped VMs)
- `create_vm` - Create a virtual machine with specified container disk (supports OS name lookup), optional instancetype and preference
- `list_instancetypes` - List available instance types
- `get_vm_instancetype` - Get instance type for a VM

#### Container Disk Lookup
The `create_vm` tool supports both full container image URLs and OS name shortcuts:
- **OS Names**: `"fedora"`, `"ubuntu"`, `"centos"`, `"debian"`, `"rhel"`, `"opensuse"`, `"alpine"`, `"cirros"`, `"windows"`, `"freebsd"`
- **Full URLs**: `"quay.io/containerdisks/fedora:latest"`, `"my-registry/my-image:tag"`
- **Auto-resolve**: Unknown OS names are resolved to `quay.io/containerdisks/{name}:latest`

### MCP Resources
- `kubevirt://{namespace}/vms` - JSON list of VMs with summary info
- `kubevirt://{namespace}/vm/{name}` - Complete VM specification
- `kubevirt://{namespace}/vm/{name}/status` - VM status and phase information
- `kubevirt://{namespace}/vmis` - JSON list of VMIs with runtime info
- `kubevirt://{namespace}/vmi/{name}` - Complete VMI specification
- `kubevirt://{namespace}/datavolumes` - JSON list of DataVolumes with source and storage info
- `kubevirt://{namespace}/datavolume/{name}` - Complete DataVolume specification

## Development

### Building
```bash
# Using Makefile (recommended)
make build

# Or directly with go
go build -o kubevirt-mcp-server .
```

### Available Make Targets
- `make build` - Build the binary (default)
- `make clean` - Clean build artifacts
- `make test` - Run tests with Ginkgo framework
- `make coverage` - Generate test coverage report
- `make fmt` - Format Go code
- `make vet` - Run go vet
- `make lint` - Run golangci-lint
- `make deps` - Download and tidy dependencies
- `make run` - Build and run the server
- `make check` - Run fmt, vet, lint, and test
- `make cluster-up` - Start kubevirtci cluster for testing
- `make cluster-down` - Stop kubevirtci cluster
- `make cluster-sync` - Build and run MCP server locally with kubevirtci access
- `make test-functional` - Run functional tests against MCP server
- `make help` - Show help message

### Testing
The project uses the Ginkgo testing framework with Gomega assertions:

```bash
# Run unit tests (with verbose output)
make test

# Generate coverage report
make coverage

# Run functional tests against MCP server
make test-functional
```

### Code Quality
The project uses golangci-lint for comprehensive code analysis:

```bash
# Run linter
make lint

# Run all quality checks (fmt + vet + lint + test)
make check
```

### Local Development Environment

For functional testing with a real KubeVirt cluster, use the kubevirtci integration:

```bash
# Start a local kubevirtci cluster with KubeVirt
make cluster-up

# Stop the cluster when done
make cluster-down

# Build and run MCP server locally with cluster access
make cluster-sync
```

The kubevirtci integration includes:

**`scripts/kubevirtci.sh`** handles:
- Downloading and setting up kubevirtci
- Starting a local Kubernetes cluster with KubeVirt and CDI
- Configuring the environment for testing
- Providing access to kubectl, kubeconfig, and registry

**`scripts/sync.sh`** handles:
- Building the MCP server binary
- Starting the MCP server locally with proper KUBECONFIG environment
- Process management (start/stop with PID tracking)
- Logging to /tmp/kubevirt-mcp-server.log for debugging

#### Test Structure
- `pkg/client/client_test.go` - Tests for KubeVirt client creation
- `pkg/tools/tools_test.go` - Tests for MCP tool handlers (argument validation)
- `pkg/resources/resources_test.go` - Tests for MCP resource handlers (URI parsing)
- `tests/functional/` - Functional tests for MCP server stdio communication
  - `functional_suite_test.go` - Test suite setup and KubeVirt cluster verification
  - `mcp_server_stdio_test.go` - Comprehensive functional tests for all MCP server functionality including:
    - MCP server initialization and tool listing
    - All MCP tools: list_vms, start_vm, stop_vm, restart_vm, list_instancetypes, get_vm_instancetype
    - All MCP resources: kubevirt://namespace/vms, vm/name, vmis, vmi/name endpoints
    - Error handling for invalid tools, missing arguments, invalid URIs, and non-existent VMs

## Git Commit Guidelines

### Conventional Commits Standard

This project follows the [Conventional Commits v1.0.0](https://www.conventionalcommits.org/en/v1.0.0/) specification for structured commit messages:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

#### Commit Types
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code (white-space, formatting, etc)
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `test`: Adding missing tests or correcting existing tests
- `build`: Changes that affect the build system or external dependencies
- `ci`: Changes to CI configuration files and scripts
- `chore`: Other changes that don't modify src or test files
- `revert`: Reverts a previous commit

#### Examples
```
feat(tools): add VM restart functionality
fix(resources): correct URI parsing for empty namespaces
docs: update README with new MCP resource endpoints
test: add missing test cases for client package
refactor(pkg): reorganize internal package structure
```

### AI Assistance Attribution

When creating commits that include AI assistance, use the `Assisted-By` trailer in commit messages:

```
feat(resources): add VM instance filtering

Add support for filtering VMs by status and labels in resource handlers.

Assisted-By: Claude <noreply@anthropic.com>
```

This provides proper attribution for AI-assisted contributions while maintaining commit history transparency.

## Documentation Maintenance

### CRITICAL REQUIREMENT
When making ANY changes to the project (code, structure, features, tools, resources), you MUST update BOTH documentation files:

1. **CLAUDE.md** - Internal development context and guidelines
2. **README.md** - Public-facing project documentation

### What to Update
- **Project Structure** - Any new files, directories, or reorganization
- **Features** - New tools, resources, or capabilities added
- **Development** - Build instructions, testing procedures, dependencies
- **Architecture** - Changes to package organization or design patterns

### When to Update
- Adding new MCP tools or resources
- Refactoring code structure
- Changing build/test procedures
- Adding new dependencies
- Modifying API endpoints or protocols

This ensures documentation stays current and accurate for both developers and users.

## Claude CLI Integration

### Configuration

The MCP server can be configured with Claude CLI in several ways:

1. **Project-specific** - `.clauderc` file in project directory
2. **Global** - `~/.config/claude/config.json` 
3. **Environment variables** - Dynamic configuration

### Usage Patterns

The server enables these workflow patterns with Claude CLI:

- **VM lifecycle management** - Start/stop/restart VMs with natural language
- **Development environment control** - Manage dev/staging VM environments  
- **Troubleshooting assistance** - Investigate VM issues with Claude's help
- **Bulk operations** - Mass VM management across namespaces
- **Resource analysis** - Review VM configurations and instance types

### Security Model

- Uses KUBECONFIG credentials for KubeVirt API access
- Inherits permissions from the configured kubeconfig
- Supports dedicated service accounts for restricted access
- No cluster admin privileges required - only VM operations

This integration enables natural language VM management workflows while maintaining security through standard Kubernetes RBAC.

# important-instruction-reminders
Do what has been asked; nothing more, nothing less.
NEVER create files unless they're absolutely necessary for achieving your goal.
ALWAYS prefer editing an existing file to creating a new one.
NEVER proactively create documentation files (*.md) or README files. Only create documentation files if explicitly requested by the User.
ALWAYS update both CLAUDE.md and README.md when making project changes that affect structure, features, or usage.
