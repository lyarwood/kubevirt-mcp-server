package tools

import (
	"context"
	"fmt"
	"github.com/lyarwood/kubevirt-mcp-server/pkg/client"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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

func VmRestart(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	// Get the VM to check its current state
	vm, err := virtClient.VirtualMachine(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return newToolResultErr(err)
	}

	// Check if VM has a running VMI
	_, err = virtClient.VirtualMachineInstance(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		// If VMI doesn't exist, just start the VM
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
					Text: fmt.Sprintf("started %s (was not running)", name),
				},
			},
		}, nil
	}

	// If VMI exists, restart by deleting the VMI (VM will recreate it)
	err = virtClient.VirtualMachineInstance(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return newToolResultErr(err)
	}

	// Ensure VM is set to restart by setting RunStrategy to Always
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
				Text: fmt.Sprintf("restarted %s", name),
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

// ResolveContainerDisk resolves OS names to container disk images from quay.io/containerdisks
func ResolveContainerDisk(input string) string {
	// If input already looks like a container image, return as-is
	if strings.Contains(input, "/") || strings.Contains(input, ":") {
		return input
	}

	// Common OS name mappings to containerdisk images
	osMap := map[string]string{
		"fedora":   "quay.io/containerdisks/fedora:latest",
		"ubuntu":   "quay.io/containerdisks/ubuntu:latest",
		"centos":   "quay.io/containerdisks/centos:latest",
		"debian":   "quay.io/containerdisks/debian:latest",
		"rhel":     "quay.io/containerdisks/rhel:latest",
		"opensuse": "quay.io/containerdisks/opensuse:latest",
		"alpine":   "quay.io/containerdisks/alpine:latest",
		"cirros":   "quay.io/kubevirt/cirros-container-disk-demo",
		"windows":  "quay.io/containerdisks/windows:latest",
		"freebsd":  "quay.io/containerdisks/freebsd:latest",
	}

	// Normalize input to lowercase for lookup
	normalized := strings.ToLower(strings.TrimSpace(input))

	// Look up the OS name
	if containerDisk, exists := osMap[normalized]; exists {
		return containerDisk
	}

	// If no match found, assume it's already a valid container disk name
	// and try to construct a containerdisks URL
	return fmt.Sprintf("quay.io/containerdisks/%s:latest", normalized)
}

func VmCreate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
	cd := request.Params.Arguments["container_disk"]
	containerDiskInput, ok := cd.(string)
	if !ok {
		return newToolResultErr(fmt.Errorf("unable to decode container_disk string"))
	}

	// Resolve the container disk image (handles OS names like "fedora", "ubuntu", etc.)
	containerDisk := ResolveContainerDisk(containerDiskInput)

	vm := &virtv1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: virtv1.VirtualMachineSpec{
			RunStrategy: &[]virtv1.VirtualMachineRunStrategy{virtv1.RunStrategyHalted}[0],
			Template: &virtv1.VirtualMachineInstanceTemplateSpec{
				Spec: virtv1.VirtualMachineInstanceSpec{
					Domain: virtv1.DomainSpec{
						Devices: virtv1.Devices{
							Disks: []virtv1.Disk{
								{
									Name: "containerdisk",
									DiskDevice: virtv1.DiskDevice{
										Disk: &virtv1.DiskTarget{
											Bus: "virtio",
										},
									},
								},
							},
						},
					},
					Volumes: []virtv1.Volume{
						{
							Name: "containerdisk",
							VolumeSource: virtv1.VolumeSource{
								ContainerDisk: &virtv1.ContainerDiskSource{
									Image: containerDisk,
								},
							},
						},
					},
				},
			},
		},
	}

	// Only set memory resources if no instancetype is provided
	// Instancetypes define their own resource requirements
	hasInstancetype := false

	if it := request.Params.Arguments["instancetype"]; it != nil {
		instancetype, ok := it.(string)
		if !ok {
			return newToolResultErr(fmt.Errorf("unable to decode instancetype string"))
		}
		vm.Spec.Instancetype = &virtv1.InstancetypeMatcher{
			Name: instancetype,
			Kind: "VirtualMachineClusterInstancetype",
		}
		hasInstancetype = true
	}

	if pref := request.Params.Arguments["preference"]; pref != nil {
		preference, ok := pref.(string)
		if !ok {
			return newToolResultErr(fmt.Errorf("unable to decode preference string"))
		}
		vm.Spec.Preference = &virtv1.PreferenceMatcher{
			Name: preference,
			Kind: "VirtualMachineClusterPreference",
		}
	}

	// Set default memory only if no instancetype is provided
	if !hasInstancetype {
		vm.Spec.Template.Spec.Domain.Resources = virtv1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse("128Mi"),
			},
		}
	}

	_, err = virtClient.VirtualMachine(namespace).Create(ctx, vm, metav1.CreateOptions{})
	if err != nil {
		return newToolResultErr(err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("created VM %s in namespace %s", name, namespace),
			},
		},
	}, nil
}
