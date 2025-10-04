package vm

import (
	"context"
	"fmt"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type StartInput struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type StartOutput struct {
	Result string `json:"result"`
}

func Start(ctx context.Context, req *mcp.CallToolRequest, input StartInput) (*mcp.CallToolResult, *StartOutput, error) {
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

	patchData := []byte(`[{"op": "replace", "path": "/spec/runStrategy", "value": "Always"}]`)
	_, err = virtClient.VirtualMachine(input.Namespace).Patch(ctx, input.Name, types.JSONPatchType, patchData, metav1.PatchOptions{})
	if err != nil {
		return nil, nil, err
	}

	return nil, &StartOutput{Result: fmt.Sprintf("started %s", input.Name)}, nil
}