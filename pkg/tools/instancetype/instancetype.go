package instancetype

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ListInput struct{}

type ListOutput struct {
	Names string `json:"names"`
}

func List(ctx context.Context, req *mcp.CallToolRequest, input ListInput) (*mcp.CallToolResult, *ListOutput, error) {
	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return nil, nil, err
	}

	instancetypes, err := virtClient.VirtualMachineClusterInstancetype().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, nil, err
	}

	var names string
	for _, instancetype := range instancetypes.Items {
		names += fmt.Sprintf("%s\n", instancetype.Name)
	}

	return nil, &ListOutput{Names: names}, nil
}

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

	instancetype, err := virtClient.VirtualMachineClusterInstancetype().Get(ctx, input.Name, metav1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}

	result := map[string]interface{}{
		"name":        instancetype.Name,
		"labels":      instancetype.Labels,
		"annotations": instancetype.Annotations,
		"spec":        instancetype.Spec,
	}

	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, nil, err
	}

	return nil, &GetOutput{Result: string(resultJSON)}, nil
}