# KubeVirt MCP server

This is a simple MCP server implementation for the KubeVirt project, a virtualization extension for Kubernetes.

## Project Structure

- `main.go` - MCP server setup and registration
- `pkg/client/` - Shared KubeVirt client utilities
- `pkg/tools/` - MCP tool handlers for VM operations
- `pkg/resources/` - MCP resource handlers for structured data access
- `Makefile` - Build automation and development tasks
- `server_config.json` - MCP server configuration for client integration
- `go.mod` / `go.sum` - Go module dependencies

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
- `make help` - Show help message

### Testing
The project uses the Ginkgo testing framework with Gomega assertions:

```bash
# Run all tests (with verbose output)
make test

# Generate coverage report
make coverage
```

### Code Quality
The project uses golangci-lint for comprehensive code analysis:

```bash
# Run linter
make lint

# Run all quality checks (fmt + vet + lint + test)
make check
```

#### Test Structure
- `pkg/client/client_test.go` - Tests for KubeVirt client creation
- `pkg/tools/tools_test.go` - Tests for MCP tool handlers (argument validation)
- `pkg/resources/resources_test.go` - Tests for MCP resource handlers (URI parsing)

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

🤖 Generated with [Claude Code](https://claude.ai/code)

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

# important-instruction-reminders
Do what has been asked; nothing more, nothing less.
NEVER create files unless they're absolutely necessary for achieving your goal.
ALWAYS prefer editing an existing file to creating a new one.
NEVER proactively create documentation files (*.md) or README files. Only create documentation files if explicitly requested by the User.
ALWAYS update both CLAUDE.md and README.md when making project changes that affect structure, features, or usage.
