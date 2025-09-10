package vm

import (
	"context"
	"fmt"
	"github.com/lyarwood/kubevirt-mcp-server/pkg/client"
	"github.com/lyarwood/kubevirt-mcp-server/pkg/tools/container"

	"github.com/mark3labs/mcp-go/mcp"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
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

func List(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return newToolResultErr(err)
	}

	namespace, err := request.RequireString("namespace")
	if err != nil {
		return newToolResultErr(fmt.Errorf("namespace parameter required: %w", err))
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

func Start(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	// Use JSON patch to update RunStrategy to avoid conflicts
	patchData := []byte(`[{"op": "replace", "path": "/spec/runStrategy", "value": "Always"}]`)
	_, err = virtClient.VirtualMachine(namespace).Patch(ctx, name, types.JSONPatchType, patchData, metav1.PatchOptions{})
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

func Stop(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	// Use JSON patch to update RunStrategy to avoid conflicts
	patchData := []byte(`[{"op": "replace", "path": "/spec/runStrategy", "value": "Halted"}]`)
	_, err = virtClient.VirtualMachine(namespace).Patch(ctx, name, types.JSONPatchType, patchData, metav1.PatchOptions{})
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

func Restart(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	// Check if VM has a running VMI
	_, err = virtClient.VirtualMachineInstance(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		// If VMI doesn't exist, just start the VM
		// Use JSON patch to update RunStrategy to avoid conflicts
		patchData := []byte(`[{"op": "replace", "path": "/spec/runStrategy", "value": "Always"}]`)
		_, err = virtClient.VirtualMachine(namespace).Patch(ctx, name, types.JSONPatchType, patchData, metav1.PatchOptions{})
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
	// Use JSON patch to update RunStrategy to avoid conflicts
	patchData := []byte(`[{"op": "replace", "path": "/spec/runStrategy", "value": "Always"}]`)
	_, err = virtClient.VirtualMachine(namespace).Patch(ctx, name, types.JSONPatchType, patchData, metav1.PatchOptions{})
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

func GetInstancetype(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

func Create(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
	containerDiskInput, err := request.RequireString("container_disk")
	if err != nil {
		return newToolResultErr(fmt.Errorf("container_disk parameter required: %w", err))
	}

	// Resolve the container disk image (handles OS names like "fedora", "ubuntu", etc.)
	containerDisk := container.ResolveContainerDisk(containerDiskInput)

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

	instancetype := request.GetString("instancetype", "")
	if instancetype != "" {
		vm.Spec.Instancetype = &virtv1.InstancetypeMatcher{
			Name: instancetype,
			Kind: "VirtualMachineClusterInstancetype",
		}
		hasInstancetype = true
	}

	preference := request.GetString("preference", "")
	if preference != "" {
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