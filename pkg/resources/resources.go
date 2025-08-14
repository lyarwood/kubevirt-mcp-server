package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"main/pkg/client"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func VmsList(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	// Parse namespace from URI: kubevirt://{namespace}/vms
	parts := strings.Split(request.Params.URI, "/")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid URI format, expected kubevirt://{namespace}/vms")
	}
	namespace := parts[2]

	virtClient, err := client.GetKubevirtClient()
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

func VmGet(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	// Parse namespace and name from URI: kubevirt://{namespace}/vm/{name}
	parts := strings.Split(request.Params.URI, "/")
	if len(parts) < 5 {
		return nil, fmt.Errorf("invalid URI format, expected kubevirt://{namespace}/vm/{name}")
	}
	namespace := parts[2]
	name := parts[4]

	virtClient, err := client.GetKubevirtClient()
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

func VmisList(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	// Parse namespace from URI: kubevirt://{namespace}/vmis
	parts := strings.Split(request.Params.URI, "/")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid URI format, expected kubevirt://{namespace}/vmis")
	}
	namespace := parts[2]

	virtClient, err := client.GetKubevirtClient()
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

func VmiGet(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	// Parse namespace and name from URI: kubevirt://{namespace}/vmi/{name}
	parts := strings.Split(request.Params.URI, "/")
	if len(parts) < 5 {
		return nil, fmt.Errorf("invalid URI format, expected kubevirt://{namespace}/vmi/{name}")
	}
	namespace := parts[2]
	name := parts[4]

	virtClient, err := client.GetKubevirtClient()
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
