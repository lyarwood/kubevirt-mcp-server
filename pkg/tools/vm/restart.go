package vm

import (
	"context"
	"fmt"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type RestartInput struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type RestartOutput struct {
	Result string `json:"result"`
}

func Restart(ctx context.Context, req *mcp.CallToolRequest, input RestartInput) (*mcp.CallToolResult, *RestartOutput, error) {
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

	_, err = virtClient.VirtualMachineInstance(input.Namespace).Get(ctx, input.Name, metav1.GetOptions{})
	if err != nil {
		patchData := []byte(`[{"op": "replace", "path": "/spec/runStrategy", "value": "Always"}]`)
		_, err = virtClient.VirtualMachine(input.Namespace).Patch(ctx, input.Name, types.JSONPatchType, patchData, metav1.PatchOptions{})
		if err != nil {
			return nil, nil, err
		}
		return nil, &RestartOutput{Result: fmt.Sprintf("started %s (was not running)", input.Name)}, nil
	}

	err = virtClient.VirtualMachineInstance(input.Namespace).Delete(ctx, input.Name, metav1.DeleteOptions{})
	if err != nil {
		return nil, nil, err
	}

	patchData := []byte(`[{"op": "replace", "path": "/spec/runStrategy", "value": "Always"}]`)
	_, err = virtClient.VirtualMachine(input.Namespace).Patch(ctx, input.Name, types.JSONPatchType, patchData, metav1.PatchOptions{})
	if err != nil {
		return nil, nil, err
	}

	return nil, &RestartOutput{Result: fmt.Sprintf("restarted %s", input.Name)}, nil
}