package vm

import (
	"context"
	"fmt"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	virtv1 "kubevirt.io/api/core/v1"
)

type UnpauseInput struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type UnpauseOutput struct {
	Result string `json:"result"`
}

func Unpause(ctx context.Context, req *mcp.CallToolRequest, input UnpauseInput) (*mcp.CallToolResult, *UnpauseOutput, error) {
	if input.Namespace == "" {
		return nil, nil, fmt.Errorf("namespace parameter is required")
	}
	if input.Name == "" {
		return nil, nil, fmt.Errorf("name parameter is required")
	}

	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return nil, nil, err
	}

	vmi, err := virtClient.VirtualMachineInstance(input.Namespace).Get(ctx, input.Name, metav1.GetOptions{})
	if err == nil && vmi != nil {
		err = virtClient.VirtualMachineInstance(input.Namespace).Unpause(ctx, vmi.Name, &virtv1.UnpauseOptions{})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to unpause VMI: %w", err)
		}
	}

	patchData := []byte(`[{"op": "replace", "path": "/spec/runStrategy", "value": "Always"}]`)
	_, err = virtClient.VirtualMachine(input.Namespace).Patch(ctx, input.Name, types.JSONPatchType, patchData, metav1.PatchOptions{})
	if err != nil {
		return nil, nil, err
	}

	return nil, &UnpauseOutput{
		Result: fmt.Sprintf("unpaused VM %s in namespace %s", input.Name, input.Namespace),
	}, nil
}