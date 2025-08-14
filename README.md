# kubevirt-mcp-server

A simple Model Context Protocol server for KubeVirt.

## Architecture

The project is organized into modular packages:

- `main.go` - MCP server setup and registration
- `pkg/client/` - Shared KubeVirt client utilities
- `pkg/tools/` - MCP tool handlers for VM operations
- `pkg/resources/` - MCP resource handlers for structured data access

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
go build -o kubevirt-mcp-server .
```

## Demo

This short demo uses mcp-cli as a bridge between the kubevirt-mcp-server and LLM.

The model used by the demo is llama3.2 running locally under ollama.

![demo](demo.gif)

## Links 

- https://www.anthropic.com/news/model-context-protocol
- https://github.com/mark3labs/mcp-go
- https://github.com/chrishayuk/mcp-cli
