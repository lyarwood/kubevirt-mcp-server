package prompts

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TroubleshootVM(ctx context.Context, request *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	namespace, ok := request.Params.Arguments["namespace"]
	if !ok || namespace == "" {
		return nil, fmt.Errorf("namespace parameter is required")
	}
	name, ok := request.Params.Arguments["name"]
	if !ok || name == "" {
		return nil, fmt.Errorf("name parameter is required")
	}

	// Optional issue description parameter
	issueDescription, _ := request.Params.Arguments["issue_description"]

	description := fmt.Sprintf("Comprehensive troubleshooting analysis for virtual machine %s in namespace %s", name, namespace)

	var prompt string
	if issueDescription != "" {
		prompt = fmt.Sprintf(`Perform comprehensive troubleshooting analysis for virtual machine %s in namespace %s.

**Reported Issue:** %s

## Diagnostic Analysis

### 1. Current Status Assessment
- Check VM phase and ready state using get_vm_phase
- Analyze all condition statuses with get_vm_conditions
- Review detailed status information using get_vm_status
- Identify any failing conditions with reasons and messages

### 2. Configuration Validation
- Verify instance type compatibility using get_vm_instancetype and get_instancetype
- Check network interface configuration for migration support
- Validate resource allocation appropriateness
- Review security and access policies

### 3. Resource Analysis
- Assess CPU and memory allocation from instance type specifications
- Check for resource constraints or mismatches
- Analyze migration capabilities (LiveMigratable, StorageLiveMigratable)
- Evaluate instance type suitability for workload

### 4. Health Indicators
- Guest agent connectivity status
- Network connectivity validation
- Storage accessibility and configuration
- Operating system compatibility with preferences

## Issue Identification
Based on the comprehensive analysis, identify:
- Any failing conditions with detailed reasons and messages
- Configuration mismatches or conflicts
- Resource allocation problems or bottlenecks
- Network interface or storage access issues
- Instance type suitability concerns

## Root Cause Analysis
Correlate findings to determine the root cause:
- Configuration vs. operational issues
- Resource limitations vs. configuration problems
- Timing issues vs. persistent failures
- Infrastructure vs. application layer problems

## Actionable Recommendations
Provide specific, prioritized recommendations:
- **Immediate Actions**: Critical fixes needed now
- **Configuration Changes**: Patches or adjustments required
- **Resource Adjustments**: Scaling or reallocation needed
- **Preventive Measures**: Steps to avoid future issues

## Quick Resolution Steps
If appropriate, suggest immediate actions:
- VM restart requirements and procedures
- Configuration patches using patch_vm tool
- Resource scaling recommendations
- Network or storage troubleshooting steps

Use all available MCP tools systematically to gather comprehensive diagnostic information and provide expert-level troubleshooting guidance with specific, actionable solutions.`, name, namespace, issueDescription)
	} else {
		prompt = fmt.Sprintf(`Perform comprehensive troubleshooting analysis for virtual machine %s in namespace %s.

## Diagnostic Analysis

### 1. Current Status Assessment
- Check VM phase and ready state using get_vm_phase
- Analyze all condition statuses with get_vm_conditions
- Review detailed status information using get_vm_status
- Identify any failing conditions with reasons and messages

### 2. Configuration Validation
- Verify instance type compatibility using get_vm_instancetype and get_instancetype
- Check network interface configuration for migration support
- Validate resource allocation appropriateness
- Review security and access policies

### 3. Resource Analysis
- Assess CPU and memory allocation from instance type specifications
- Check for resource constraints or mismatches
- Analyze migration capabilities (LiveMigratable, StorageLiveMigratable)
- Evaluate instance type suitability for workload

### 4. Health Indicators
- Guest agent connectivity status
- Network connectivity validation
- Storage accessibility and configuration
- Operating system compatibility with preferences

## Issue Identification
Based on the comprehensive analysis, identify:
- Any failing conditions with detailed reasons and messages
- Configuration mismatches or conflicts
- Resource allocation problems or bottlenecks
- Network interface or storage access issues
- Instance type suitability concerns

## Root Cause Analysis
Correlate findings to determine the root cause:
- Configuration vs. operational issues
- Resource limitations vs. configuration problems
- Timing issues vs. persistent failures
- Infrastructure vs. application layer problems

## Actionable Recommendations
Provide specific, prioritized recommendations:
- **Immediate Actions**: Critical fixes needed now
- **Configuration Changes**: Patches or adjustments required
- **Resource Adjustments**: Scaling or reallocation needed
- **Preventive Measures**: Steps to avoid future issues

## Quick Resolution Steps
If appropriate, suggest immediate actions:
- VM restart requirements and procedures
- Configuration patches using patch_vm tool
- Resource scaling recommendations
- Network or storage troubleshooting steps

Use all available MCP tools systematically to gather comprehensive diagnostic information and provide expert-level troubleshooting guidance with specific, actionable solutions.`, name, namespace)
	}

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