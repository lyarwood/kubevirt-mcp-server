package main

import (
	"context"
	"fmt"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/prompts"
	"github.com/lyarwood/kubevirt-mcp-server/pkg/resources"
	"github.com/lyarwood/kubevirt-mcp-server/pkg/tools/instancetype"
	"github.com/lyarwood/kubevirt-mcp-server/pkg/tools/preference"
	"github.com/lyarwood/kubevirt-mcp-server/pkg/tools/vm"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// Create MCP server
	s := mcp.NewServer(
		&mcp.Implementation{
			Name:    "kubevirt-mcp-server",
			Version: "0.0.1",
		},
		&mcp.ServerOptions{
			HasResources: true,
			HasPrompts:   true,
			HasTools:     true,
		},
	)

	mcp.AddTool(s, &mcp.Tool{Name: "list_vms", Description: "list the names of virtual machine within a given namespace"}, vm.List)
	mcp.AddTool(s, &mcp.Tool{Name: "start_vm", Description: "start the virtual machine with a given name in the provided namespace"}, vm.Start)
	mcp.AddTool(s, &mcp.Tool{Name: "stop_vm", Description: "stop the virtual machine with a given name in the provided namespace"}, vm.Stop)
	mcp.AddTool(s, &mcp.Tool{Name: "restart_vm", Description: "restart the virtual machine with a given name in the provided namespace"}, vm.Restart)
	mcp.AddTool(s, &mcp.Tool{Name: "list_instancetypes", Description: "list the name of all instance types"}, instancetype.List)
	mcp.AddTool(s, &mcp.Tool{Name: "get_instancetype", Description: "get detailed information about a specific instance type including CPU, memory, annotations and labels"}, instancetype.Get)
	mcp.AddTool(s, &mcp.Tool{Name: "get_preference", Description: "get detailed information about a specific preference including settings, annotations and labels"}, preference.Get)
	mcp.AddTool(s, &mcp.Tool{Name: "get_vm_instancetype", Description: "show the name of the instance type referenced by a virtual machine"}, vm.GetInstancetype)
	mcp.AddTool(s, &mcp.Tool{Name: "get_vm_status", Description: "get comprehensive status information for a virtual machine including ready state, generation, and state change requests"}, vm.GetStatus)
	mcp.AddTool(s, &mcp.Tool{Name: "get_vm_conditions", Description: "get detailed condition information for a virtual machine including health checks and operational state"}, vm.GetConditions)
	mcp.AddTool(s, &mcp.Tool{Name: "get_vm_phase", Description: "get current phase and basic status information for a virtual machine"}, vm.GetPhase)
	mcp.AddTool(s, &mcp.Tool{Name: "patch_vm", Description: "apply a JSON merge patch to modify a virtual machine configuration"}, vm.Patch)
	mcp.AddTool(s, &mcp.Tool{Name: "create_vm", Description: "create a virtual machine with the given name, container disk image (supports OS names like 'fedora', 'ubuntu'), and optional instancetype and preference"}, vm.Create)
	mcp.AddTool(s, &mcp.Tool{Name: "delete_vm", Description: "delete the virtual machine with a given name in the provided namespace"}, vm.Delete)
	mcp.AddTool(s, &mcp.Tool{Name: "pause_vm", Description: "pause the virtual machine with a given name in the provided namespace"}, vm.Pause)
	mcp.AddTool(s, &mcp.Tool{Name: "unpause_vm", Description: "unpause the virtual machine with a given name in the provided namespace"}, vm.Unpause)
	mcp.AddTool(s, &mcp.Tool{Name: "get_vm_disks", Description: "get the list of disks for a specified virtual machine in the given namespace"}, vm.Disks)

	// Add MCP Resource Templates
	s.AddResourceTemplate(&mcp.ResourceTemplate{URITemplate: "kubevirt://{namespace}/vms", Title: "Virtual Machines", Description: "List of virtual machines in a namespace", MIMEType: "application/json"}, resources.VmsList)
	s.AddResourceTemplate(&mcp.ResourceTemplate{URITemplate: "kubevirt://{namespace}/vm/{name}", Title: "Virtual Machine", Description: "Individual virtual machine details", MIMEType: "application/json"}, resources.VmGet)
	s.AddResourceTemplate(&mcp.ResourceTemplate{URITemplate: "kubevirt://{namespace}/vmis", Title: "Virtual Machine Instances", Description: "List of virtual machine instances in a namespace", MIMEType: "application/json"}, resources.VmisList)
	s.AddResourceTemplate(&mcp.ResourceTemplate{URITemplate: "kubevirt://{namespace}/vmi/{name}", Title: "Virtual Machine Instance", Description: "Individual virtual machine instance details", MIMEType: "application/json"}, resources.VmiGet)
	s.AddResourceTemplate(&mcp.ResourceTemplate{URITemplate: "kubevirt://{namespace}/datavolumes", Title: "Data Volumes", Description: "List of data volumes with source and storage information", MIMEType: "application/json"}, resources.DataVolumesList)
	s.AddResourceTemplate(&mcp.ResourceTemplate{URITemplate: "kubevirt://{namespace}/datavolume/{name}", Title: "Data Volume", Description: "Individual data volume specification", MIMEType: "application/json"}, resources.DataVolumeGet)
	s.AddResourceTemplate(&mcp.ResourceTemplate{URITemplate: "kubevirt://{namespace}/vm/{name}/status", Title: "VM Status", Description: "Virtual machine status and phase information", MIMEType: "application/json"}, resources.VmGetStatus)
	s.AddResourceTemplate(&mcp.ResourceTemplate{URITemplate: "kubevirt://{namespace}/vmi/{name}/guestosinfo", Title: "VMI Guest OS Info", Description: "Virtual machine instance guest operating system information", MIMEType: "application/json"}, resources.VmiGetGuestOSInfo)
	s.AddResourceTemplate(&mcp.ResourceTemplate{URITemplate: "kubevirt://{namespace}/vmi/{name}/filesystems", Title: "VMI Filesystems", Description: "Virtual machine instance filesystem information", MIMEType: "application/json"}, resources.VmiGetFilesystems)
	s.AddResourceTemplate(&mcp.ResourceTemplate{URITemplate: "kubevirt://{namespace}/vmi/{name}/userlist", Title: "VMI User List", Description: "Virtual machine instance user list information", MIMEType: "application/json"}, resources.VmiGetUserList)
	s.AddResourceTemplate(&mcp.ResourceTemplate{URITemplate: "kubevirt://{namespace}/vm/{name}/console", Title: "VM Console", Description: "Virtual machine console connection details", MIMEType: "application/json"}, resources.VmGetConsole)
	s.AddResourceTemplate(&mcp.ResourceTemplate{URITemplate: "kubevirt://{namespace}/instancetypes", Title: "Instance Types", Description: "List of instance types in a namespace", MIMEType: "application/json"}, resources.InstancetypesList)
	s.AddResourceTemplate(&mcp.ResourceTemplate{URITemplate: "kubevirt://{namespace}/preferences", Title: "Preferences", Description: "List of VM preferences in a namespace", MIMEType: "application/json"}, resources.PreferencesList)
	s.AddResourceTemplate(&mcp.ResourceTemplate{URITemplate: "kubevirt://cluster/instancetypes", Title: "Cluster Instance Types", Description: "List of cluster-wide instance types", MIMEType: "application/json"}, resources.ClusterInstancetypesList)
	s.AddResourceTemplate(&mcp.ResourceTemplate{URITemplate: "kubevirt://cluster/preferences", Title: "Cluster Preferences", Description: "List of cluster-wide VM preferences", MIMEType: "application/json"}, resources.ClusterPreferencesList)
	s.AddResourceTemplate(&mcp.ResourceTemplate{URITemplate: "kubevirt://cluster/instancetype/{name}", Title: "Cluster Instance Type", Description: "Individual cluster instance type specification", MIMEType: "application/json"}, resources.ClusterInstancetypeGet)
	s.AddResourceTemplate(&mcp.ResourceTemplate{URITemplate: "kubevirt://cluster/preference/{name}", Title: "Cluster Preference", Description: "Individual cluster preference specification", MIMEType: "application/json"}, resources.ClusterPreferenceGet)

	// Add MCP Prompts
	s.AddPrompt(&mcp.Prompt{
		Name:        "describe_vm",
		Description: "Provide a comprehensive description of a virtual machine including its configuration, status, and operational details",
		Arguments: []*mcp.PromptArgument{
			{Name: "namespace", Description: "The namespace containing the virtual machine", Required: true},
			{Name: "name", Description: "The name of the virtual machine to describe", Required: true},
		},
	}, prompts.DescribeVM)

	s.AddPrompt(&mcp.Prompt{
		Name:        "troubleshoot_vm",
		Description: "Diagnose and analyze potential issues with a virtual machine, providing actionable recommendations",
		Arguments: []*mcp.PromptArgument{
			{Name: "namespace", Description: "The namespace containing the virtual machine", Required: true},
			{Name: "name", Description: "The name of the virtual machine to troubleshoot", Required: true},
			{Name: "issue_description", Description: "Optional description of the specific issue being experienced"},
		},
	}, prompts.TroubleshootVM)

	s.AddPrompt(&mcp.Prompt{
		Name:        "health_check_vm",
		Description: "Perform a quick health check on a virtual machine and report any issues",
		Arguments: []*mcp.PromptArgument{
			{Name: "namespace", Description: "The namespace containing the virtual machine", Required: true},
			{Name: "name", Description: "The name of the virtual machine to check", Required: true},
		},
	}, prompts.HealthCheckVM)

	// Start the stdio server
	if err := s.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}