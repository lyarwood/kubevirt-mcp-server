package main

import (
	"fmt"
	
	"github.com/lyarwood/kubevirt-mcp-server/pkg/resources"
	"github.com/lyarwood/kubevirt-mcp-server/pkg/tools"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
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
		tools.VmsList,
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
		tools.VmStart,
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
		tools.VmStop,
	)

	s.AddTool(
		mcp.NewTool(
			"list_instancetypes",
			mcp.WithDescription("list the name of all instance types"),
		),
		tools.InstancetypesList,
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
		tools.VmGetInstancetype,
	)

	// Add MCP Resources
	s.AddResource(
		mcp.NewResource(
			"kubevirt://*/vms",
			"Virtual Machines",
			mcp.WithResourceDescription("List of virtual machines in a namespace"),
			mcp.WithMIMEType("application/json"),
		),
		resources.VmsList,
	)

	s.AddResource(
		mcp.NewResource(
			"kubevirt://*/vm/*",
			"Virtual Machine",
			mcp.WithResourceDescription("Individual virtual machine details"),
			mcp.WithMIMEType("application/json"),
		),
		resources.VmGet,
	)

	s.AddResource(
		mcp.NewResource(
			"kubevirt://*/vmis",
			"Virtual Machine Instances",
			mcp.WithResourceDescription("List of virtual machine instances in a namespace"),
			mcp.WithMIMEType("application/json"),
		),
		resources.VmisList,
	)

	s.AddResource(
		mcp.NewResource(
			"kubevirt://*/vmi/*",
			"Virtual Machine Instance",
			mcp.WithResourceDescription("Individual virtual machine instance details"),
			mcp.WithMIMEType("application/json"),
		),
		resources.VmiGet,
	)

	// TODO prompt
	// describe virtual machine ?

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
