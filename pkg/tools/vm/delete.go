package vm

import (
	"context"
	"fmt"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeleteInput struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type DeleteOutput struct {
	Result string `json:"result"`
}

func Delete(ctx context.Context, req *mcp.CallToolRequest, input DeleteInput) (*mcp.CallToolResult, *DeleteOutput, error) {
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

	err = virtClient.VirtualMachine(input.Namespace).Delete(ctx, input.Name, metav1.DeleteOptions{})
	if err != nil {
		return nil, nil, err
	}

	return nil, &DeleteOutput{
		Result: fmt.Sprintf("deleted VM %s in namespace %s", input.Name, input.Namespace),
	}, nil
}