package vm

import (
	"context"
	"fmt"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/client"
	"github.com/mark3labs/mcp-go/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func Restart(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	// Check if VM has a running VMI
	_, err = virtClient.VirtualMachineInstance(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		// If VMI doesn't exist, just start the VM
		// Use JSON patch to update RunStrategy to avoid conflicts
		patchData := []byte(`[{"op": "replace", "path": "/spec/runStrategy", "value": "Always"}]`)
		_, err = virtClient.VirtualMachine(namespace).Patch(ctx, name, types.JSONPatchType, patchData, metav1.PatchOptions{})
		if err != nil {
			return newToolResultErr(err)
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("started %s (was not running)", name),
				},
			},
		}, nil
	}

	// If VMI exists, restart by deleting the VMI (VM will recreate it)
	err = virtClient.VirtualMachineInstance(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return newToolResultErr(err)
	}

	// Ensure VM is set to restart by setting RunStrategy to Always
	// Use JSON patch to update RunStrategy to avoid conflicts
	patchData := []byte(`[{"op": "replace", "path": "/spec/runStrategy", "value": "Always"}]`)
	_, err = virtClient.VirtualMachine(namespace).Patch(ctx, name, types.JSONPatchType, patchData, metav1.PatchOptions{})
	if err != nil {
		return newToolResultErr(err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("restarted %s", name),
			},
		},
	}, nil
}
