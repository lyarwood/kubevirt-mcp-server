package prompts

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

func HealthCheckVM(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	namespace, ok := request.Params.Arguments["namespace"]
	if !ok || namespace == "" {
		return nil, fmt.Errorf("namespace parameter is required")
	}
	name, ok := request.Params.Arguments["name"]
	if !ok || name == "" {
		return nil, fmt.Errorf("name parameter is required")
	}

	description := fmt.Sprintf("Quick health assessment of virtual machine %s in namespace %s", name, namespace)

	prompt := fmt.Sprintf(`Perform a rapid health assessment of VM %s in namespace %s using available MCP tools.

## Quick Health Check Status Report

### 🔍 Status Indicators (use get_vm_status and get_vm_phase)
- [ ] VM Ready Status (check if ready: true)
- [ ] VM Running Status (check if status: "Running")  
- [ ] Generation Sync (desiredGeneration == observedGeneration)
- [ ] No Pending State Changes (check stateChangeRequests)

### 🔍 Condition Health (use get_vm_conditions)
- [ ] Ready Condition (status: "True")
- [ ] Guest Agent Connected (status: "True")
- [ ] Live Migration Capable (status: "True" - indicates healthy networking)
- [ ] Storage Migration Ready (status: "True")
- [ ] No Failed Conditions (all conditions should be True or have acceptable reasons)

### 🔍 Configuration Health (use get_vm_instancetype and get_instancetype)
- [ ] Valid Instance Type Assignment
- [ ] Appropriate Resource Allocation
- [ ] Compatible Instance Type Configuration
- [ ] Proper CPU/Memory Balance

### 🔍 Operational Readiness
- [ ] Network Configuration Supports Migration
- [ ] No Recent Configuration Conflicts
- [ ] Stable Resource Allocation
- [ ] Guest OS Compatibility

## Health Summary Format
After checking all indicators, provide a clear status:

**Overall Health**: ✅ HEALTHY / ⚠️ WARNING / ❌ CRITICAL

**Key Findings**:
- List 2-3 most important status points
- Highlight any immediate concerns
- Note any configuration recommendations

**Immediate Actions** (if needed):
- Specific steps to resolve critical issues
- Configuration adjustments required
- Monitoring recommendations

**Migration Readiness**: ✅ READY / ⚠️ LIMITED / ❌ BLOCKED
- Brief explanation of migration capabilities

Focus on rapid assessment with clear pass/fail indicators and immediate actionable insights. Use the MCP tools efficiently to gather essential health information without extensive analysis.`, name, namespace)

	messages := []mcp.PromptMessage{
		mcp.NewPromptMessage(mcp.RoleUser, mcp.TextContent{
			Type: "text",
			Text: prompt,
		}),
	}

	return mcp.NewGetPromptResult(description, messages), nil
}