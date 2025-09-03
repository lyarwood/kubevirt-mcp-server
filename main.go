package main

import (
	"flag"
	"log"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/resources"
	"github.com/lyarwood/kubevirt-mcp-server/pkg/tools"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Parse command line flags
	var (
		httpAddr = flag.String("http", "", "HTTP server address (e.g., ':8080'). If empty, uses stdio transport.")
	)
	flag.Parse()

	// Create MCP server
	s := server.NewMCPServer(
		"kubevirt MCP server demo 🚀",
		"0.0.1",
		server.WithResourceCapabilities(true, true),
		server.WithToolCapabilities(true),
		server.WithLogging(),
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
			"restart_vm",
			mcp.WithDescription("restart the virtual machine with a given name in the provided namespace"),
			mcp.WithString(
				"namespace",
				mcp.Description("The namespace of the virtual machine"),
				mcp.Required()),
			mcp.WithString(
				"name",
				mcp.Description("The name of the virtual machine"),
				mcp.Required()),
		),
		tools.VmRestart,
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

	s.AddTool(
		mcp.NewTool(
			"create_vm",
			mcp.WithDescription("create a virtual machine with the given name, container disk image (supports OS names like 'fedora', 'ubuntu'), and optional instancetype and preference"),
			mcp.WithString(
				"namespace",
				mcp.Description("The namespace for the virtual machine"),
				mcp.Required()),
			mcp.WithString(
				"name",
				mcp.Description("The name of the virtual machine"),
				mcp.Required()),
			mcp.WithString(
				"container_disk",
				mcp.Description("The container disk image to use for the VM (supports OS names like 'fedora', 'ubuntu' or full URLs)"),
				mcp.Required()),
			mcp.WithString(
				"instancetype",
				mcp.Description("Optional instance type name")),
			mcp.WithString(
				"preference",
				mcp.Description("Optional preference name")),
		),
		tools.VmCreate,
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

	s.AddResource(
		mcp.NewResource(
			"kubevirt://*/datavolumes",
			"Data Volumes",
			mcp.WithResourceDescription("List of data volumes with source and storage information"),
			mcp.WithMIMEType("application/json"),
		),
		resources.DataVolumesList,
	)

	s.AddResource(
		mcp.NewResource(
			"kubevirt://*/datavolume/*",
			"Data Volume",
			mcp.WithResourceDescription("Individual data volume specification"),
			mcp.WithMIMEType("application/json"),
		),
		resources.DataVolumeGet,
	)

	s.AddResource(
		mcp.NewResource(
			"kubevirt://*/vm/*/status",
			"VM Status",
			mcp.WithResourceDescription("Virtual machine status and phase information"),
			mcp.WithMIMEType("application/json"),
		),
		resources.VmGetStatus,
	)

	// TODO prompt
	// describe virtual machine ?

	// Start the appropriate server based on command line flags
	if *httpAddr != "" {
		// Start HTTP server
		log.Printf("Starting KubeVirt MCP HTTP server on %s", *httpAddr)
		httpServer := server.NewStreamableHTTPServer(s, 
			server.WithStateLess(true), // Enable stateless mode for easier testing
		)
		if err := httpServer.Start(*httpAddr); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	} else {
		// Start stdio server (default)
		log.Printf("Starting KubeVirt MCP stdio server")
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("Stdio server error: %v", err)
		}
	}
}
