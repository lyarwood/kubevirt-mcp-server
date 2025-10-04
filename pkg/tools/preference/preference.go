package preference

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GetInput struct {
	Name string `json:"name"`
}

type GetOutput struct {
	Result string `json:"result"`
}

func Get(ctx context.Context, req *mcp.CallToolRequest, input GetInput) (*mcp.CallToolResult, *GetOutput, error) {
	if input.Name == "" {
		return nil, nil, fmt.Errorf("name parameter is required")
	}

	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return nil, nil, err
	}

	preference, err := virtClient.VirtualMachineClusterPreference().Get(ctx, input.Name, metav1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}

	result := map[string]interface{}{
		"name":        preference.Name,
		"labels":      preference.Labels,
		"annotations": preference.Annotations,
		"spec":        preference.Spec,
	}

	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, nil, err
	}

	return nil, &GetOutput{Result: string(resultJSON)}, nil
}