package prompts

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func DescribeVM(ctx context.Context, request *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	namespace, ok := request.Params.Arguments["namespace"]
	if !ok || namespace == "" {
		return nil, fmt.Errorf("namespace parameter is required")
	}
	name, ok := request.Params.Arguments["name"]
	if !ok || name == "" {
		return nil, fmt.Errorf("name parameter is required")
	}

	description := fmt.Sprintf("Comprehensive description of virtual machine %s in namespace %s", name, namespace)

	prompt := fmt.Sprintf(`Analyze the virtual machine %s in namespace %s and provide a comprehensive description including:

## VM Overview
- Current status and phase using get_vm_phase and get_vm_status tools
- Instance type details using get_vm_instancetype and get_instancetype tools  
- Operating system preferences and configuration
- Creation time and generation information

## Configuration Details
- Network interfaces and connectivity type
- Resource allocation (CPU, memory) from instance type
- Live migration capabilities and requirements
- Security and access settings

## Operational Status
- Health conditions and readiness using get_vm_conditions
- Migration capabilities (LiveMigratable, StorageLiveMigratable)
- Recent state changes or pending operations
- Guest agent connectivity status

## Resource Analysis
- Instance type characteristics and use case suitability
- Network configuration impact on migration
- Storage configuration and accessibility
- Performance and scaling considerations

Please use the available MCP tools (get_vm_status, get_vm_conditions, get_vm_phase, get_vm_instancetype, get_instancetype) to gather comprehensive information and present it in a clear, organized format suitable for both technical review and operational planning.

Focus on providing actionable insights about the VM's current state, configuration optimizations, and operational readiness.`, name, namespace)

	messages := []*mcp.PromptMessage{
		{
			Role: "user",
			Content: &mcp.TextContent{
				Text: prompt,
			},
		},
	}

	return &mcp.GetPromptResult{
		Description: description,
		Messages:    messages,
	}, nil
}