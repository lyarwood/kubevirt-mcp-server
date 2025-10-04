package vm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type PatchInput struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Patch     string `json:"patch"`
}

type PatchOutput struct {
	Result string `json:"result"`
}

func Patch(ctx context.Context, req *mcp.CallToolRequest, input PatchInput) (*mcp.CallToolResult, *PatchOutput, error) {
	if input.Namespace == "" {
		return nil, nil, fmt.Errorf("namespace parameter is required")
	}
	if input.Name == "" {
		return nil, nil, fmt.Errorf("name parameter is required")
	}
	if input.Patch == "" {
		return nil, nil, fmt.Errorf("patch parameter is required")
	}

	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return nil, nil, err
	}

	// Validate that patch is valid JSON
	var patchJSON interface{}
	if err := json.Unmarshal([]byte(input.Patch), &patchJSON); err != nil {
		return nil, nil, fmt.Errorf("invalid JSON in patch parameter: %w", err)
	}

	// Get the current VM to validate it exists
	currentVM, err := virtClient.VirtualMachine(input.Namespace).Get(ctx, input.Name, metav1.GetOptions{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get VM %s/%s: %w", input.Namespace, input.Name, err)
	}

	// Apply the patch
	patchedVM, err := virtClient.VirtualMachine(input.Namespace).Patch(ctx, input.Name, types.MergePatchType, []byte(input.Patch), metav1.PatchOptions{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to patch VM %s/%s: %w", input.Namespace, input.Name, err)
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
		return nil, nil, err
	}

	return nil, &PatchOutput{Result: string(resultJSON)}, nil
}