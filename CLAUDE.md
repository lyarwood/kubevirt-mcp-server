# KubeVirt MCP server

This is a simple MCP server implementation for the KubeVirt project, a virtualization extension for Kubernetes.

## Project Structure

- `main.go` - MCP server setup and registration
- `pkg/client/` - Shared KubeVirt client utilities
- `pkg/tools/` - MCP tool handlers for VM operations
- `pkg/resources/` - MCP resource handlers for structured data access
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
go build -o kubevirt-mcp-server .
```

### Testing
Build the project to verify syntax and dependencies are correct.

## Git Commit Guidelines

When creating commits that include AI assistance, use the `Assisted-By` trailer in commit messages:

```
Add new feature implementation

Description of changes made.

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
