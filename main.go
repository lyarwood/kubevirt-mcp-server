package main

import (
	"fmt"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/resources"
	"github.com/lyarwood/kubevirt-mcp-server/pkg/tools/instancetype"
	"github.com/lyarwood/kubevirt-mcp-server/pkg/tools/preference"
	"github.com/lyarwood/kubevirt-mcp-server/pkg/tools/vm"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"kubevirt MCP server demo 🚀",
		"0.0.1",
		server.WithResourceCapabilities(true, true),
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
		vm.List,
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
		vm.Start,
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
		vm.Stop,
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
		vm.Restart,
	)

	s.AddTool(
		mcp.NewTool(
			"list_instancetypes",
			mcp.WithDescription("list the name of all instance types"),
		),
		instancetype.List,
	)

	s.AddTool(
		mcp.NewTool(
			"get_instancetype",
			mcp.WithDescription("get detailed information about a specific instance type including CPU, memory, annotations and labels"),
			mcp.WithString(
				"name",
				mcp.Description("The name of the instance type"),
				mcp.Required()),
		),
		instancetype.Get,
	)

	s.AddTool(
		mcp.NewTool(
			"get_preference",
			mcp.WithDescription("get detailed information about a specific preference including settings, annotations and labels"),
			mcp.WithString(
				"name",
				mcp.Description("The name of the preference"),
				mcp.Required()),
		),
		preference.Get,
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
		vm.GetInstancetype,
	)

	s.AddTool(
		mcp.NewTool(
			"get_vm_status",
			mcp.WithDescription("get comprehensive status information for a virtual machine including ready state, generation, and state change requests"),
			mcp.WithString(
				"namespace",
				mcp.Description("The namespace of the virtual machine"),
				mcp.Required()),
			mcp.WithString(
				"name",
				mcp.Description("The name of the virtual machine"),
				mcp.Required()),
		),
		vm.GetStatus,
	)

	s.AddTool(
		mcp.NewTool(
			"get_vm_conditions",
			mcp.WithDescription("get detailed condition information for a virtual machine including health checks and operational state"),
			mcp.WithString(
				"namespace",
				mcp.Description("The namespace of the virtual machine"),
				mcp.Required()),
			mcp.WithString(
				"name",
				mcp.Description("The name of the virtual machine"),
				mcp.Required()),
		),
		vm.GetConditions,
	)

	s.AddTool(
		mcp.NewTool(
			"get_vm_phase",
			mcp.WithDescription("get current phase and basic status information for a virtual machine"),
			mcp.WithString(
				"namespace",
				mcp.Description("The namespace of the virtual machine"),
				mcp.Required()),
			mcp.WithString(
				"name",
				mcp.Description("The name of the virtual machine"),
				mcp.Required()),
		),
		vm.GetPhase,
	)

	s.AddTool(
		mcp.NewTool(
			"patch_vm",
			mcp.WithDescription("apply a JSON merge patch to modify a virtual machine configuration"),
			mcp.WithString(
				"namespace",
				mcp.Description("The namespace of the virtual machine"),
				mcp.Required()),
			mcp.WithString(
				"name",
				mcp.Description("The name of the virtual machine"),
				mcp.Required()),
			mcp.WithString(
				"patch",
				mcp.Description("JSON merge patch data to apply to the VM (e.g., network interface modifications, resource updates)"),
				mcp.Required()),
		),
		vm.Patch,
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
		vm.Create,
	)

	s.AddTool(
		mcp.NewTool(
			"delete_vm",
			mcp.WithDescription("delete the virtual machine with a given name in the provided namespace"),
			mcp.WithString(
				"namespace",
				mcp.Description("The namespace of the virtual machine"),
				mcp.Required()),
			mcp.WithString(
				"name",
				mcp.Description("The name of the virtual machine"),
				mcp.Required()),
		),
		vm.Delete,
	)

	s.AddTool(
		mcp.NewTool(
			"pause_vm",
			mcp.WithDescription("pause the virtual machine with a given name in the provided namespace"),
			mcp.WithString(
				"namespace",
				mcp.Description("The namespace of the virtual machine"),
				mcp.Required()),
			mcp.WithString(
				"name",
				mcp.Description("The name of the virtual machine"),
				mcp.Required()),
		),
		vm.Pause,
	)

	s.AddTool(
		mcp.NewTool(
			"unpause_vm",
			mcp.WithDescription("unpause the virtual machine with a given name in the provided namespace"),
			mcp.WithString(
				"namespace",
				mcp.Description("The namespace of the virtual machine"),
				mcp.Required()),
			mcp.WithString(
				"name",
				mcp.Description("The name of the virtual machine"),
				mcp.Required()),
		),
		vm.Unpause,
	)

	// Add MCP Resource Templates
	s.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"kubevirt://{namespace}/vms",
			"Virtual Machines",
			mcp.WithTemplateDescription("List of virtual machines in a namespace"),
			mcp.WithTemplateMIMEType("application/json"),
		),
		resources.VmsList,
	)

	s.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"kubevirt://{namespace}/vm/{name}",
			"Virtual Machine",
			mcp.WithTemplateDescription("Individual virtual machine details"),
			mcp.WithTemplateMIMEType("application/json"),
		),
		resources.VmGet,
	)

	s.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"kubevirt://{namespace}/vmis",
			"Virtual Machine Instances",
			mcp.WithTemplateDescription("List of virtual machine instances in a namespace"),
			mcp.WithTemplateMIMEType("application/json"),
		),
		resources.VmisList,
	)

	s.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"kubevirt://{namespace}/vmi/{name}",
			"Virtual Machine Instance",
			mcp.WithTemplateDescription("Individual virtual machine instance details"),
			mcp.WithTemplateMIMEType("application/json"),
		),
		resources.VmiGet,
	)

	s.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"kubevirt://{namespace}/datavolumes",
			"Data Volumes",
			mcp.WithTemplateDescription("List of data volumes with source and storage information"),
			mcp.WithTemplateMIMEType("application/json"),
		),
		resources.DataVolumesList,
	)

	s.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"kubevirt://{namespace}/datavolume/{name}",
			"Data Volume",
			mcp.WithTemplateDescription("Individual data volume specification"),
			mcp.WithTemplateMIMEType("application/json"),
		),
		resources.DataVolumeGet,
	)

	s.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"kubevirt://{namespace}/vm/{name}/status",
			"VM Status",
			mcp.WithTemplateDescription("Virtual machine status and phase information"),
			mcp.WithTemplateMIMEType("application/json"),
		),
		resources.VmGetStatus,
	)

	s.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"kubevirt://{namespace}/vmi/{name}/guestosinfo",
			"VMI Guest OS Info",
			mcp.WithTemplateDescription("Virtual machine instance guest operating system information"),
			mcp.WithTemplateMIMEType("application/json"),
		),
		resources.VmiGetGuestOSInfo,
	)

	s.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"kubevirt://{namespace}/vmi/{name}/filesystems",
			"VMI Filesystems",
			mcp.WithTemplateDescription("Virtual machine instance filesystem information"),
			mcp.WithTemplateMIMEType("application/json"),
		),
		resources.VmiGetFilesystems,
	)

	s.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"kubevirt://{namespace}/vmi/{name}/userlist",
			"VMI User List",
			mcp.WithTemplateDescription("Virtual machine instance user list information"),
			mcp.WithTemplateMIMEType("application/json"),
		),
		resources.VmiGetUserList,
	)

	s.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"kubevirt://{namespace}/vm/{name}/console",
			"VM Console",
			mcp.WithTemplateDescription("Virtual machine console connection details"),
			mcp.WithTemplateMIMEType("application/json"),
		),
		resources.VmGetConsole,
	)

	s.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"kubevirt://{namespace}/instancetypes",
			"Instance Types",
			mcp.WithTemplateDescription("List of instance types in a namespace"),
			mcp.WithTemplateMIMEType("application/json"),
		),
		resources.InstancetypesList,
	)

	s.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"kubevirt://{namespace}/preferences",
			"Preferences",
			mcp.WithTemplateDescription("List of VM preferences in a namespace"),
			mcp.WithTemplateMIMEType("application/json"),
		),
		resources.PreferencesList,
	)

	s.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"kubevirt://cluster/instancetypes",
			"Cluster Instance Types",
			mcp.WithTemplateDescription("List of cluster-wide instance types"),
			mcp.WithTemplateMIMEType("application/json"),
		),
		resources.ClusterInstancetypesList,
	)

	s.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"kubevirt://cluster/preferences",
			"Cluster Preferences",
			mcp.WithTemplateDescription("List of cluster-wide VM preferences"),
			mcp.WithTemplateMIMEType("application/json"),
		),
		resources.ClusterPreferencesList,
	)

	s.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"kubevirt://cluster/instancetype/{name}",
			"Cluster Instance Type",
			mcp.WithTemplateDescription("Individual cluster instance type specification"),
			mcp.WithTemplateMIMEType("application/json"),
		),
		resources.ClusterInstancetypeGet,
	)

	s.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"kubevirt://cluster/preference/{name}",
			"Cluster Preference",
			mcp.WithTemplateDescription("Individual cluster preference specification"),
			mcp.WithTemplateMIMEType("application/json"),
		),
		resources.ClusterPreferenceGet,
	)

	// TODO prompt
	// describe virtual machine ?

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
