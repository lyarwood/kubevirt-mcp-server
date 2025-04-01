package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	virtv1 "kubevirt.io/api/core/v1"
	"kubevirt.io/client-go/kubecli"
)

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"kubevirt MCP server demo 🚀",
		"0.0.1",
	)

	// TODO resources
	// kubevirt://{namespace}/vms
	// kubevirt://{namespace}/vm/{name}
	// kubevirt://{namespace}/vmis
	// kubevirt://{namespace}/vmi/{name}

	// TODO tools
	// list instance types
	// change instance type
	s.AddTool(
		mcp.NewTool(
			"list_vm_names",
			mcp.WithDescription("list the names of virtual machine within a given namespace"),
			mcp.WithString(
				"namespace",
				mcp.Description("The namespace of the virtual machine"),
				mcp.Required()),
		),
		vmsListTool,
	)

	s.AddTool(
		mcp.NewTool(
			"start_vm",
			mcp.WithDescription("start the virtual machine with a given name in the provided namespace"),
			mcp.WithString(
				"namespace",
				mcp.Description("The namespace of the virtual machine"),
				mcp.Required()),
			mcp.WithString(
				"name",
				mcp.Description("The Name of the virtual machine"),
				mcp.Required()),
		),
		vmStartTool,
	)

	s.AddTool(
		mcp.NewTool(
			"stop_vm",
			mcp.WithDescription("stop the virtual machine with a given name in the provided namespace"),
			mcp.WithString(
				"namespace",
				mcp.Description("The namespace of the virtual machine"),
				mcp.Required()),
			mcp.WithString(
				"name",
				mcp.Description("The name of the virtual machine"),
				mcp.Required()),
		),
		vmStopTool,
	)

	s.AddTool(
		mcp.NewTool(
			"list_instancetypes",
			mcp.WithDescription("list the name of all instance types"),
		),
		instancetypeListTool,
	)

	s.AddTool(
		mcp.NewTool(
			"get_vm_instancetype",
			mcp.WithDescription("show the name of the instance type referenced by a virtual machine"),
			mcp.WithString(
				"namespace",
				mcp.Description("The namespace of the virtual machine"),
				mcp.Required()),
			mcp.WithString(
				"name",
				mcp.Description("The name of the virtual machine"),
				mcp.Required()),
		),
		vmGetInstancetype,
	)

	// TODO prompt
	// describe virtual machine ?

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func vmGetInstancetype(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clientConfig := kubecli.DefaultClientConfig(&pflag.FlagSet{})
	virtClient, err := kubecli.GetKubevirtClientFromClientConfig(clientConfig)
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

	message := "no instasnce type referenced by virtual machine"
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

func instancetypeListTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clientConfig := kubecli.DefaultClientConfig(&pflag.FlagSet{})
	virtClient, err := kubecli.GetKubevirtClientFromClientConfig(clientConfig)
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

func vmStartTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clientConfig := kubecli.DefaultClientConfig(&pflag.FlagSet{})
	virtClient, err := kubecli.GetKubevirtClientFromClientConfig(clientConfig)
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

func vmStopTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clientConfig := kubecli.DefaultClientConfig(&pflag.FlagSet{})
	virtClient, err := kubecli.GetKubevirtClientFromClientConfig(clientConfig)
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

func vmsListTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clientConfig := kubecli.DefaultClientConfig(&pflag.FlagSet{})
	virtClient, err := kubecli.GetKubevirtClientFromClientConfig(clientConfig)
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
