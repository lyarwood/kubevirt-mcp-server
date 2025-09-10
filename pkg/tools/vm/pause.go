package vm

import (
	"context"
	"fmt"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/client"
	"github.com/mark3labs/mcp-go/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	virtv1 "kubevirt.io/api/core/v1"
)

func Pause(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	// Use JSON patch to update RunStrategy to Manual and set paused state
	patchData := []byte(`[{"op": "replace", "path": "/spec/runStrategy", "value": "Manual"}, {"op": "replace", "path": "/spec/running", "value": false}]`)
	_, err = virtClient.VirtualMachine(namespace).Patch(ctx, name, types.JSONPatchType, patchData, metav1.PatchOptions{})
	if err != nil {
		return newToolResultErr(err)
	}

	// Now pause the VMI if it exists
	vmi, err := virtClient.VirtualMachineInstance(namespace).Get(ctx, name, metav1.GetOptions{})
	if err == nil && vmi != nil {
		// Pause the VMI using subresource
		err = virtClient.VirtualMachineInstance(namespace).Pause(ctx, vmi.Name, &virtv1.PauseOptions{})
		if err != nil {
			return newToolResultErr(fmt.Errorf("failed to pause VMI: %w", err))
		}
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("paused VM %s in namespace %s", name, namespace),
			},
		},
	}, nil
}
