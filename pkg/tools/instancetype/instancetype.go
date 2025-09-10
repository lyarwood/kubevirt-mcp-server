package instancetype

import (
	"context"
	"fmt"
	"github.com/lyarwood/kubevirt-mcp-server/pkg/client"

	"github.com/mark3labs/mcp-go/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newToolResultErr(err error) (*mcp.CallToolResult, error) {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: err.Error(),
			},
		},
	}, err
}

func List(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return newToolResultErr(err)
	}

	instancetypes, err := virtClient.VirtualMachineClusterInstancetype().List(ctx, metav1.ListOptions{})
	if err != nil {
		return newToolResultErr(err)
	}

	names := ""
	for _, instancetype := range instancetypes.Items {
		names += fmt.Sprintf("%s\n", instancetype.Name)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: names,
			},
		},
	}, nil
}