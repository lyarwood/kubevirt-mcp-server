package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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
		server.WithResourceCapabilities(true, true),
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
			"list_vms",
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

	// Add MCP Resources
	s.AddResource(
		mcp.NewResource(
			"kubevirt://*/vms",
			"Virtual Machines",
			mcp.WithResourceDescription("List of virtual machines in a namespace"),
			mcp.WithMIMEType("application/json"),
		),
		vmsResource,
	)

	s.AddResource(
		mcp.NewResource(
			"kubevirt://*/vm/*",
			"Virtual Machine",
			mcp.WithResourceDescription("Individual virtual machine details"),
			mcp.WithMIMEType("application/json"),
		),
		vmResource,
	)

	s.AddResource(
		mcp.NewResource(
			"kubevirt://*/vmis",
			"Virtual Machine Instances",
			mcp.WithResourceDescription("List of virtual machine instances in a namespace"),
			mcp.WithMIMEType("application/json"),
		),
		vmisResource,
	)

	s.AddResource(
		mcp.NewResource(
			"kubevirt://*/vmi/*",
			"Virtual Machine Instance",
			mcp.WithResourceDescription("Individual virtual machine instance details"),
			mcp.WithMIMEType("application/json"),
		),
		vmiResource,
	)

	// TODO prompt
	// describe virtual machine ?

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

// Resource handlers

func vmsResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	// Parse namespace from URI: kubevirt://{namespace}/vms
	parts := strings.Split(request.Params.URI, "/")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid URI format, expected kubevirt://{namespace}/vms")
	}
	namespace := parts[2]

	clientConfig := kubecli.DefaultClientConfig(&pflag.FlagSet{})
	virtClient, err := kubecli.GetKubevirtClientFromClientConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	vms, err := virtClient.VirtualMachine(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	vmList := make([]map[string]interface{}, 0, len(vms.Items))
	for _, vm := range vms.Items {
		vmInfo := map[string]interface{}{
			"name":      vm.Name,
			"namespace": vm.Namespace,
			"status":    vm.Status.PrintableStatus,
			"created":   vm.CreationTimestamp,
		}
		if vm.Spec.RunStrategy != nil {
			vmInfo["runStrategy"] = string(*vm.Spec.RunStrategy)
		}
		if vm.Spec.Instancetype != nil {
			vmInfo["instanceType"] = vm.Spec.Instancetype.Name
		}
		vmList = append(vmList, vmInfo)
	}

	jsonData, err := json.MarshalIndent(vmList, "", "  ")
	if err != nil {
		return nil, err
	}

	return []mcp.ResourceContents{
		&mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(jsonData),
		},
	}, nil
}

func vmResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	// Parse namespace and name from URI: kubevirt://{namespace}/vm/{name}
	parts := strings.Split(request.Params.URI, "/")
	if len(parts) < 5 {
		return nil, fmt.Errorf("invalid URI format, expected kubevirt://{namespace}/vm/{name}")
	}
	namespace := parts[2]
	name := parts[4]

	clientConfig := kubecli.DefaultClientConfig(&pflag.FlagSet{})
	virtClient, err := kubecli.GetKubevirtClientFromClientConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	vm, err := virtClient.VirtualMachine(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	jsonData, err := json.MarshalIndent(vm, "", "  ")
	if err != nil {
		return nil, err
	}

	return []mcp.ResourceContents{
		&mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(jsonData),
		},
	}, nil
}

func vmisResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	// Parse namespace from URI: kubevirt://{namespace}/vmis
	parts := strings.Split(request.Params.URI, "/")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid URI format, expected kubevirt://{namespace}/vmis")
	}
	namespace := parts[2]

	clientConfig := kubecli.DefaultClientConfig(&pflag.FlagSet{})
	virtClient, err := kubecli.GetKubevirtClientFromClientConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	vmis, err := virtClient.VirtualMachineInstance(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	vmiList := make([]map[string]interface{}, 0, len(vmis.Items))
	for _, vmi := range vmis.Items {
		vmiInfo := map[string]interface{}{
			"name":      vmi.Name,
			"namespace": vmi.Namespace,
			"phase":     vmi.Status.Phase,
			"created":   vmi.CreationTimestamp,
			"nodeName":  vmi.Status.NodeName,
		}
		if len(vmi.Status.Interfaces) > 0 {
			interfaces := make([]map[string]interface{}, 0, len(vmi.Status.Interfaces))
			for _, iface := range vmi.Status.Interfaces {
				ifaceInfo := map[string]interface{}{
					"name": iface.Name,
				}
				if iface.IP != "" {
					ifaceInfo["ip"] = iface.IP
				}
				if iface.MAC != "" {
					ifaceInfo["mac"] = iface.MAC
				}
				interfaces = append(interfaces, ifaceInfo)
			}
			vmiInfo["interfaces"] = interfaces
		}
		vmiList = append(vmiList, vmiInfo)
	}

	jsonData, err := json.MarshalIndent(vmiList, "", "  ")
	if err != nil {
		return nil, err
	}

	return []mcp.ResourceContents{
		&mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(jsonData),
		},
	}, nil
}

func vmiResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	// Parse namespace and name from URI: kubevirt://{namespace}/vmi/{name}
	parts := strings.Split(request.Params.URI, "/")
	if len(parts) < 5 {
		return nil, fmt.Errorf("invalid URI format, expected kubevirt://{namespace}/vmi/{name}")
	}
	namespace := parts[2]
	name := parts[4]

	clientConfig := kubecli.DefaultClientConfig(&pflag.FlagSet{})
	virtClient, err := kubecli.GetKubevirtClientFromClientConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	vmi, err := virtClient.VirtualMachineInstance(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	jsonData, err := json.MarshalIndent(vmi, "", "  ")
	if err != nil {
		return nil, err
	}

	return []mcp.ResourceContents{
		&mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(jsonData),
		},
	}, nil
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
