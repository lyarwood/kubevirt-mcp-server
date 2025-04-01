# kubevirt-mcp-server

A simple Model Context Protocol server for KubeVirt.

## Tools

The following tools are currently provided:

### list_vms

### start_vm

### stop_vm

### list_instancetypes

### get_vm_instancetype

## Demo

This short demo uses mcp-cli as a bridge between the kubevirt-mcp-server and LLM.

The model used by the demo is llama3.2 running locally under ollama.

![demo](demo.gif)

## Links 

- https://www.anthropic.com/news/model-context-protocol
- https://github.com/mark3labs/mcp-go
- https://github.com/chrishayuk/mcp-cli
