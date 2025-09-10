package vm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/client"
	"github.com/mark3labs/mcp-go/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func Patch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return newToolResultErr(err)
	}

	namespace, err := request.RequireString("namespace")
	if err != nil {
		return newToolResultErr(fmt.Errorf("namespace parameter required: %w", err))
	}
	name, err := request.RequireString("name")
	if err != nil {
		return newToolResultErr(fmt.Errorf("name parameter required: %w", err))
	}
	patchData, err := request.RequireString("patch")
	if err != nil {
		return newToolResultErr(fmt.Errorf("patch parameter required: %w", err))
	}

	// Validate that patch is valid JSON
	var patchJSON interface{}
	if err := json.Unmarshal([]byte(patchData), &patchJSON); err != nil {
		return newToolResultErr(fmt.Errorf("invalid JSON in patch parameter: %w", err))
	}

	// Get the current VM to validate it exists
	currentVM, err := virtClient.VirtualMachine(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return newToolResultErr(fmt.Errorf("failed to get VM %s/%s: %w", namespace, name, err))
	}

	// Apply the patch
	patchedVM, err := virtClient.VirtualMachine(namespace).Patch(ctx, name, types.MergePatchType, []byte(patchData), metav1.PatchOptions{})
	if err != nil {
		return newToolResultErr(fmt.Errorf("failed to patch VM %s/%s: %w", namespace, name, err))
	}

	// Create success response with information about what was changed
	result := map[string]interface{}{
		"name":      patchedVM.Name,
		"namespace": patchedVM.Namespace,
		"message":   "VM successfully patched",
		"generation": map[string]interface{}{
			"before": currentVM.Generation,
			"after":  patchedVM.Generation,
		},
		"resourceVersion": map[string]interface{}{
			"before": currentVM.ResourceVersion,
			"after":  patchedVM.ResourceVersion,
		},
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