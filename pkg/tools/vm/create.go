package vm

import (
	"context"
	"fmt"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/client"
	"github.com/lyarwood/kubevirt-mcp-server/pkg/tools/containerdisks"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	virtv1 "kubevirt.io/api/core/v1"
)

type CreateInput struct {
	Namespace     string `json:"namespace"`
	Name          string `json:"name"`
	ContainerDisk string `json:"container_disk"`
	Instancetype  string `json:"instancetype,omitempty"`
	Preference    string `json:"preference,omitempty"`
}

type CreateOutput struct {
	Result string `json:"result"`
}

func Create(ctx context.Context, req *mcp.CallToolRequest, input CreateInput) (*mcp.CallToolResult, *CreateOutput, error) {
	if input.Namespace == "" {
		return nil, nil, fmt.Errorf("namespace parameter is required")
	}
	if input.Name == "" {
		return nil, nil, fmt.Errorf("name parameter is required")
	}
	if input.ContainerDisk == "" {
		return nil, nil, fmt.Errorf("container_disk parameter is required")
	}

	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return nil, nil, err
	}

	containerDisk := containerdisks.ResolveContainerDisk(input.ContainerDisk)

	vm := &virtv1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      input.Name,
			Namespace: input.Namespace,
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

	hasInstancetype := false
	if input.Instancetype != "" {
		vm.Spec.Instancetype = &virtv1.InstancetypeMatcher{
			Name: input.Instancetype,
			Kind: "VirtualMachineClusterInstancetype",
		}
		hasInstancetype = true
	}

	if input.Preference != "" {
		vm.Spec.Preference = &virtv1.PreferenceMatcher{
			Name: input.Preference,
			Kind: "VirtualMachineClusterPreference",
		}
	}

	if !hasInstancetype {
		vm.Spec.Template.Spec.Domain.Resources = virtv1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse("128Mi"),
			},
		}
	}

	_, err = virtClient.VirtualMachine(input.Namespace).Create(ctx, vm, metav1.CreateOptions{})
	if err != nil {
		return nil, nil, err
	}

	return nil, &CreateOutput{
		Result: fmt.Sprintf("created VM %s in namespace %s", input.Name, input.Namespace),
	}, nil
}