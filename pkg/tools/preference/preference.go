package preference

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lyarwood/kubevirt-mcp-server/pkg/client"

	"github.com/mark3labs/mcp-go/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newToolResultErr(err error) (*mcp.CallToolResult, error) {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: err.Error(),
			},
		},
	}, err
}

func Get(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return newToolResultErr(fmt.Errorf("name parameter is required: %w", err))
	}

	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return newToolResultErr(err)
	}

	preference, err := virtClient.VirtualMachineClusterPreference().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return newToolResultErr(err)
	}

	result := map[string]interface{}{
		"name":        preference.Name,
		"labels":      preference.Labels,
		"annotations": preference.Annotations,
		"spec":        preference.Spec,
	}

	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return newToolResultErr(err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(resultJSON),
			},
		},
	}, nil
}
