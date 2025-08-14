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
