package vm

import (
	"context"
	"fmt"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ListInput struct {
	Namespace string `json:"namespace"`
}

type ListOutput struct {
	Names string `json:"names"`
}

func List(ctx context.Context, req *mcp.CallToolRequest, input ListInput) (*mcp.CallToolResult, *ListOutput, error) {
	if input.Namespace == "" {
		return nil, nil, fmt.Errorf("namespace parameter is required")
	}

	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return nil, nil, err
	}

	vms, err := virtClient.VirtualMachine(input.Namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, nil, err
	}

	var names string
	for _, vm := range vms.Items {
		names += fmt.Sprintf("%s\n", vm.Name)
	}

	return nil, &ListOutput{Names: names}, nil
}