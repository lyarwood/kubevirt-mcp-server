package vm

import (
	"context"
	"fmt"
	"strings"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/client"
	"github.com/mark3labs/mcp-go/mcp"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Disks(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return newToolResultErr(err)
	}

	namespace, err := request.RequireString("namespace")
	if err != nil {
		return newToolResultErr(fmt.Errorf("namespace parameter required: %w", err))
	}
	name, err := request.RequireString("name")
	if err != nil {
		return newToolResultErr(fmt.Errorf("name parameter required: %w", err))
	}

	vm, err := virtClient.VirtualMachine(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return newToolResultErr(err)
	}

	var diskNames []string
	for _, disk := range vm.Spec.Template.Spec.Domain.Devices.Disks {
		diskNames = append(diskNames, disk.Name)
	}

	disks := "No disks found"
	if len(diskNames) > 0 {
		disks = strings.Join(diskNames, ", ")
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: disks,
			},
		},
	}, nil
}
