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

type PauseInput struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type PauseOutput struct {
	Result string `json:"result"`
}

func Pause(ctx context.Context, req *mcp.CallToolRequest, input PauseInput) (*mcp.CallToolResult, *PauseOutput, error) {
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

	patchData := []byte(`[{"op": "replace", "path": "/spec/runStrategy", "value": "Manual"}, {"op": "replace", "path": "/spec/running", "value": false}]`)
	_, err = virtClient.VirtualMachine(input.Namespace).Patch(ctx, input.Name, types.JSONPatchType, patchData, metav1.PatchOptions{})
	if err != nil {
		return nil, nil, err
	}

	vmi, err := virtClient.VirtualMachineInstance(input.Namespace).Get(ctx, input.Name, metav1.GetOptions{})
	if err == nil && vmi != nil {
		err = virtClient.VirtualMachineInstance(input.Namespace).Pause(ctx, vmi.Name, &virtv1.PauseOptions{})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to pause VMI: %w", err)
		}
	}

	return nil, &PauseOutput{
		Result: fmt.Sprintf("paused VM %s in namespace %s", input.Name, input.Namespace),
	}, nil
}