package vm

import (
	"context"
	"fmt"
	"strings"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DisksInput struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type DisksOutput struct {
	Disks string `json:"disks"`
}

func Disks(ctx context.Context, req *mcp.CallToolRequest, input DisksInput) (*mcp.CallToolResult, *DisksOutput, error) {
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

	vm, err := virtClient.VirtualMachine(input.Namespace).Get(ctx, input.Name, metav1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}

	var diskNames []string
	for _, disk := range vm.Spec.Template.Spec.Domain.Devices.Disks {
		diskNames = append(diskNames, disk.Name)
	}

	disks := "No disks found"
	if len(diskNames) > 0 {
		disks = strings.Join(diskNames, ", ")
	}

	return nil, &DisksOutput{Disks: disks}, nil
}