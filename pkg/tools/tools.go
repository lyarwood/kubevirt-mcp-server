package tools

import (
	"context"
	"fmt"
	"main/pkg/client"

	"github.com/mark3labs/mcp-go/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	virtv1 "kubevirt.io/api/core/v1"
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

func VmsList(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return newToolResultErr(err)
	}

	ns := request.Params.Arguments["namespace"]
	namespace, ok := ns.(string)
	if !ok {
		return newToolResultErr(fmt.Errorf("unable to decode namespace string"))
	}
	vms, err := virtClient.VirtualMachine(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return newToolResultErr(err)
	}

	names := ""
	for _, vm := range vms.Items {
		names += fmt.Sprintf("%s\n", vm.Name)
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

func VmStart(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return newToolResultErr(err)
	}

	ns := request.Params.Arguments["namespace"]
	namespace, ok := ns.(string)
	if !ok {
		return newToolResultErr(fmt.Errorf("unable to decode namespace string"))
	}
	n := request.Params.Arguments["name"]
	name, ok := n.(string)
	if !ok {
		return newToolResultErr(fmt.Errorf("unable to decode name string"))
	}

	vm, err := virtClient.VirtualMachine(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return newToolResultErr(err)
	}

	s := virtv1.RunStrategyAlways
	vm.Spec.RunStrategy = &s

	_, err = virtClient.VirtualMachine(namespace).Update(ctx, vm, metav1.UpdateOptions{})
	if err != nil {
		return newToolResultErr(err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("started %s", name),
			},
		},
	}, nil
}

func VmStop(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return newToolResultErr(err)
	}

	ns := request.Params.Arguments["namespace"]
	namespace, ok := ns.(string)
	if !ok {
		return newToolResultErr(fmt.Errorf("unable to decode namespace string"))
	}
	n := request.Params.Arguments["name"]
	name, ok := n.(string)
	if !ok {
		return newToolResultErr(fmt.Errorf("unable to decode name string"))
	}

	vm, err := virtClient.VirtualMachine(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return newToolResultErr(err)
	}

	s := virtv1.RunStrategyHalted
	vm.Spec.RunStrategy = &s

	_, err = virtClient.VirtualMachine(namespace).Update(ctx, vm, metav1.UpdateOptions{})
	if err != nil {
		return newToolResultErr(err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("stopped %s", name),
			},
		},
	}, nil
}

func InstancetypesList(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

func VmGetInstancetype(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return newToolResultErr(err)
	}

	ns := request.Params.Arguments["namespace"]
	namespace, ok := ns.(string)
	if !ok {
		return newToolResultErr(fmt.Errorf("unable to decode namespace string"))
	}
	n := request.Params.Arguments["name"]
	name, ok := n.(string)
	if !ok {
		return newToolResultErr(fmt.Errorf("unable to decode name string"))
	}

	vm, err := virtClient.VirtualMachine(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return newToolResultErr(err)
	}

	message := "no instance type referenced by virtual machine"
	if vm.Spec.Instancetype != nil {
		message = vm.Spec.Instancetype.Name
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: message,
			},
		},
	}, nil
}