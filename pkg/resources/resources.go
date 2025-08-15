package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lyarwood/kubevirt-mcp-server/pkg/client"
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

func DataVolumesList(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	// Parse namespace from URI: kubevirt://{namespace}/datavolumes
	parts := strings.Split(request.Params.URI, "/")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid URI format, expected kubevirt://{namespace}/datavolumes")
	}
	namespace := parts[2]

	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return nil, err
	}

	// Get the underlying clientset to access CDI resources
	clientset := virtClient.CdiClient()
	dataVolumes, err := clientset.CdiV1beta1().DataVolumes(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	dvList := make([]map[string]interface{}, 0, len(dataVolumes.Items))
	for _, dv := range dataVolumes.Items {
		dvInfo := map[string]interface{}{
			"name":      dv.Name,
			"namespace": dv.Namespace,
			"phase":     dv.Status.Phase,
			"created":   dv.CreationTimestamp,
		}

		// Add source information if available
		if dv.Spec.Source != nil {
			source := map[string]interface{}{}
			if dv.Spec.Source.HTTP != nil {
				source["type"] = "http"
				source["url"] = dv.Spec.Source.HTTP.URL
			} else if dv.Spec.Source.S3 != nil {
				source["type"] = "s3"
				source["url"] = dv.Spec.Source.S3.URL
			} else if dv.Spec.Source.Registry != nil {
				source["type"] = "registry"
				source["url"] = dv.Spec.Source.Registry.URL
			} else if dv.Spec.Source.PVC != nil {
				source["type"] = "pvc"
				source["name"] = dv.Spec.Source.PVC.Name
				source["namespace"] = dv.Spec.Source.PVC.Namespace
			} else if dv.Spec.Source.Upload != nil {
				source["type"] = "upload"
			} else if dv.Spec.Source.Blank != nil {
				source["type"] = "blank"
			}
			if len(source) > 0 {
				dvInfo["source"] = source
			}
		}

		// Add storage information
		if dv.Spec.Storage != nil {
			storage := map[string]interface{}{}
			if dv.Spec.Storage.Resources.Requests != nil {
				if storageSize, ok := dv.Spec.Storage.Resources.Requests["storage"]; ok {
					storage["size"] = storageSize.String()
				}
			}
			if dv.Spec.Storage.StorageClassName != nil {
				storage["storageClass"] = *dv.Spec.Storage.StorageClassName
			}
			if len(storage) > 0 {
				dvInfo["storage"] = storage
			}
		}

		// Add progress information if available
		if dv.Status.Progress != "" {
			dvInfo["progress"] = dv.Status.Progress
		}

		dvList = append(dvList, dvInfo)
	}

	jsonData, err := json.MarshalIndent(dvList, "", "  ")
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

func DataVolumeGet(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	// Parse namespace and name from URI: kubevirt://{namespace}/datavolume/{name}
	parts := strings.Split(request.Params.URI, "/")
	if len(parts) < 5 {
		return nil, fmt.Errorf("invalid URI format, expected kubevirt://{namespace}/datavolume/{name}")
	}
	namespace := parts[2]
	name := parts[4]

	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return nil, err
	}

	// Get the underlying clientset to access CDI resources
	clientset := virtClient.CdiClient()
	dataVolume, err := clientset.CdiV1beta1().DataVolumes(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	jsonData, err := json.MarshalIndent(dataVolume, "", "  ")
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
