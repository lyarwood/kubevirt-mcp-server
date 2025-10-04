package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func VmsList(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	parts := strings.Split(req.Params.URI, "/")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid URI format, expected kubevirt://{namespace}/vms")
	}
	namespace := parts[2]
	if namespace == "" {
		return nil, fmt.Errorf("resource namespace may not be empty")
	}

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

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(jsonData),
			},
		},
	}, nil
}

func VmGet(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	parts := strings.Split(req.Params.URI, "/")
	if len(parts) < 5 {
		return nil, fmt.Errorf("invalid URI format, expected kubevirt://{namespace}/vm/{name}")
	}
	namespace := parts[2]
	if namespace == "" {
		return nil, fmt.Errorf("resource namespace may not be empty")
	}
	name := parts[4]
	if name == "" {
		return nil, fmt.Errorf("resource name may not be empty")
	}

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

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(jsonData),
			},
		},
	}, nil
}

func VmisList(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	parts := strings.Split(req.Params.URI, "/")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid URI format, expected kubevirt://{namespace}/vmis")
	}
	namespace := parts[2]
	if namespace == "" {
		return nil, fmt.Errorf("resource namespace may not be empty")
	}

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

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(jsonData),
			},
		},
	}, nil
}

func VmiGet(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	parts := strings.Split(req.Params.URI, "/")
	if len(parts) < 5 {
		return nil, fmt.Errorf("invalid URI format, expected kubevirt://{namespace}/vmi/{name}")
	}
	namespace := parts[2]
	if namespace == "" {
		return nil, fmt.Errorf("resource namespace may not be empty")
	}
	name := parts[4]
	if name == "" {
		return nil, fmt.Errorf("resource name may not be empty")
	}

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

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(jsonData),
			},
		},
	}, nil
}

func DataVolumesList(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	parts := strings.Split(req.Params.URI, "/")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid URI format, expected kubevirt://{namespace}/datavolumes")
	}
	namespace := parts[2]
	if namespace == "" {
		return nil, fmt.Errorf("resource namespace may not be empty")
	}

	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return nil, err
	}

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

		if dv.Status.Progress != "" {
			dvInfo["progress"] = dv.Status.Progress
		}

		dvList = append(dvList, dvInfo)
	}

	jsonData, err := json.MarshalIndent(dvList, "", "  ")
	if err != nil {
		return nil, err
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(jsonData),
			},
		},
	}, nil
}

func DataVolumeGet(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	parts := strings.Split(req.Params.URI, "/")
	if len(parts) < 5 {
		return nil, fmt.Errorf("invalid URI format, expected kubevirt://{namespace}/datavolume/{name}")
	}
	namespace := parts[2]
	if namespace == "" {
		return nil, fmt.Errorf("resource namespace may not be empty")
	}
	name := parts[4]
	if name == "" {
		return nil, fmt.Errorf("resource name may not be empty")
	}

	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return nil, err
	}

	clientset := virtClient.CdiClient()
	dataVolume, err := clientset.CdiV1beta1().DataVolumes(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	jsonData, err := json.MarshalIndent(dataVolume, "", "  ")
	if err != nil {
		return nil, err
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(jsonData),
			},
		},
	}, nil
}

func VmGetStatus(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	parts := strings.Split(req.Params.URI, "/")
	if len(parts) < 6 {
		return nil, fmt.Errorf("invalid URI format, expected kubevirt://{namespace}/vm/{name}/status")
	}
	namespace := parts[2]
	if namespace == "" {
		return nil, fmt.Errorf("resource namespace may not be empty")
	}
	name := parts[4]
	if name == "" {
		return nil, fmt.Errorf("resource name may not be empty")
	}

	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return nil, err
	}

	vm, err := virtClient.VirtualMachine(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	statusInfo := map[string]interface{}{
		"name":      vm.Name,
		"namespace": vm.Namespace,
		"status":    vm.Status.PrintableStatus,
		"phase":     vm.Status.Ready,
		"created":   vm.CreationTimestamp,
		"ready":     vm.Status.Ready,
	}

	if vm.Spec.RunStrategy != nil {
		statusInfo["runStrategy"] = string(*vm.Spec.RunStrategy)
	}

	if len(vm.Status.StateChangeRequests) > 0 {
		requests := make([]map[string]interface{}, 0, len(vm.Status.StateChangeRequests))
		for _, r := range vm.Status.StateChangeRequests {
			request := map[string]interface{}{
				"action": r.Action,
			}
			if r.UID != nil {
				request["uid"] = *r.UID
			}
			requests = append(requests, request)
		}
		statusInfo["stateChangeRequests"] = requests
	}

	if len(vm.Status.Conditions) > 0 {
		conditions := make([]map[string]interface{}, 0, len(vm.Status.Conditions))
		for _, cond := range vm.Status.Conditions {
			condition := map[string]interface{}{
				"type":   cond.Type,
				"status": cond.Status,
			}
			if cond.Reason != "" {
				condition["reason"] = cond.Reason
			}
			if cond.Message != "" {
				condition["message"] = cond.Message
			}
			if !cond.LastTransitionTime.IsZero() {
				condition["lastTransitionTime"] = cond.LastTransitionTime
			}
			conditions = append(conditions, condition)
		}
		statusInfo["conditions"] = conditions
	}

	jsonData, err := json.MarshalIndent(statusInfo, "", "  ")
	if err != nil {
		return nil, err
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(jsonData),
			},
		},
	}, nil
}

func VmiGetGuestOSInfo(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	parts := strings.Split(req.Params.URI, "/")
	if len(parts) < 6 {
		return nil, fmt.Errorf("invalid URI format, expected kubevirt://{namespace}/vmi/{name}/guestosinfo")
	}
	namespace := parts[2]
	if namespace == "" {
		return nil, fmt.Errorf("resource namespace may not be empty")
	}
	name := parts[4]
	if name == "" {
		return nil, fmt.Errorf("resource name may not be empty")
	}

	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return nil, err
	}

	guestOSInfo, err := virtClient.VirtualMachineInstance(namespace).GuestOsInfo(ctx, name)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.MarshalIndent(guestOSInfo, "", "  ")
	if err != nil {
		return nil, err
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(jsonData),
			},
		},
	}, nil
}

func VmiGetFilesystems(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	parts := strings.Split(req.Params.URI, "/")
	if len(parts) < 6 {
		return nil, fmt.Errorf("invalid URI format, expected kubevirt://{namespace}/vmi/{name}/filesystems")
	}
	namespace := parts[2]
	if namespace == "" {
		return nil, fmt.Errorf("resource namespace may not be empty")
	}
	name := parts[4]
	if name == "" {
		return nil, fmt.Errorf("resource name may not be empty")
	}

	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return nil, err
	}

	filesystems, err := virtClient.VirtualMachineInstance(namespace).FilesystemList(ctx, name)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.MarshalIndent(filesystems, "", "  ")
	if err != nil {
		return nil, err
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(jsonData),
			},
		},
	}, nil
}

func VmiGetUserList(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	parts := strings.Split(req.Params.URI, "/")
	if len(parts) < 6 {
		return nil, fmt.Errorf("invalid URI format, expected kubevirt://{namespace}/vmi/{name}/userlist")
	}
	namespace := parts[2]
	if namespace == "" {
		return nil, fmt.Errorf("resource namespace may not be empty")
	}
	name := parts[4]
	if name == "" {
		return nil, fmt.Errorf("resource name may not be empty")
	}

	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return nil, err
	}

	userList, err := virtClient.VirtualMachineInstance(namespace).UserList(ctx, name)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.MarshalIndent(userList, "", "  ")
	if err != nil {
		return nil, err
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(jsonData),
			},
		},
	}, nil
}

func VmGetConsole(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	parts := strings.Split(req.Params.URI, "/")
	if len(parts) < 6 {
		return nil, fmt.Errorf("invalid URI format, expected kubevirt://{namespace}/vm/{name}/console")
	}
	namespace := parts[2]
	if namespace == "" {
		return nil, fmt.Errorf("resource namespace may not be empty")
	}
	name := parts[4]
	if name == "" {
		return nil, fmt.Errorf("resource name may not be empty")
	}

	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return nil, err
	}

	vmi, err := virtClient.VirtualMachineInstance(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	consoleInfo := map[string]interface{}{
		"name":      vmi.Name,
		"namespace": vmi.Namespace,
		"phase":     vmi.Status.Phase,
		"nodeName":  vmi.Status.NodeName,
	}

	consoles := []string{}
	if vmi.Status.Phase == "Running" {
		consoles = append(consoles, "vnc", "serial")
		for _, condition := range vmi.Status.Conditions {
			if condition.Type == "AgentConnected" && condition.Status == "True" {
				consoles = append(consoles, "guest-agent")
				break
			}
		}
	}
	consoleInfo["availableConsoles"] = consoles

	if len(consoles) > 0 {
		connectionInfo := map[string]interface{}{
			"note": "Use kubectl or virtctl to connect to consoles",
			"commands": map[string]string{
				"vnc":    fmt.Sprintf("virtctl vnc %s -n %s", name, namespace),
				"serial": fmt.Sprintf("virtctl console %s -n %s", name, namespace),
			},
		}
		consoleInfo["connectionInfo"] = connectionInfo
	}

	jsonData, err := json.MarshalIndent(consoleInfo, "", "  ")
	if err != nil {
		return nil, err
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(jsonData),
			},
		},
	}, nil
}

func InstancetypesList(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	parts := strings.Split(req.Params.URI, "/")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid URI format, expected kubevirt://{namespace}/instancetypes")
	}
	namespace := parts[2]
	if namespace == "" {
		return nil, fmt.Errorf("resource namespace may not be empty")
	}

	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return nil, err
	}

	instancetypes, err := virtClient.VirtualMachineInstancetype(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	instancetypeList := make([]map[string]interface{}, 0, len(instancetypes.Items))
	for _, it := range instancetypes.Items {
		itInfo := map[string]interface{}{
			"name":      it.Name,
			"namespace": it.Namespace,
			"created":   it.CreationTimestamp,
		}
		if it.Spec.CPU.Guest != 0 {
			itInfo["cpu"] = it.Spec.CPU.Guest
		}
		if !it.Spec.Memory.Guest.IsZero() {
			itInfo["memory"] = it.Spec.Memory.Guest.String()
		}
		instancetypeList = append(instancetypeList, itInfo)
	}

	jsonData, err := json.MarshalIndent(instancetypeList, "", "  ")
	if err != nil {
		return nil, err
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(jsonData),
			},
		},
	}, nil
}

func PreferencesList(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	parts := strings.Split(req.Params.URI, "/")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid URI format, expected kubevirt://{namespace}/preferences")
	}
	namespace := parts[2]
	if namespace == "" {
		return nil, fmt.Errorf("resource namespace may not be empty")
	}

	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return nil, err
	}

	preferences, err := virtClient.VirtualMachinePreference(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	preferenceList := make([]map[string]interface{}, 0, len(preferences.Items))
	for _, pref := range preferences.Items {
		prefInfo := map[string]interface{}{
			"name":      pref.Name,
			"namespace": pref.Namespace,
			"created":   pref.CreationTimestamp,
		}
		if pref.Spec.Machine != nil {
			prefInfo["hasMachinePreferences"] = true
		}
		preferenceList = append(preferenceList, prefInfo)
	}

	jsonData, err := json.MarshalIndent(preferenceList, "", "  ")
	if err != nil {
		return nil, err
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(jsonData),
			},
		},
	}, nil
}

func ClusterInstancetypesList(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	parts := strings.Split(req.Params.URI, "/")
	if len(parts) < 3 || parts[2] != "cluster" {
		return nil, fmt.Errorf("invalid URI format, expected kubevirt://cluster/instancetypes")
	}

	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return nil, err
	}

	instancetypes, err := virtClient.VirtualMachineClusterInstancetype().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	instancetypeList := make([]map[string]interface{}, 0, len(instancetypes.Items))
	for _, it := range instancetypes.Items {
		itInfo := map[string]interface{}{
			"name":    it.Name,
			"created": it.CreationTimestamp,
		}
		if it.Spec.CPU.Guest != 0 {
			itInfo["cpu"] = it.Spec.CPU.Guest
		}
		if !it.Spec.Memory.Guest.IsZero() {
			itInfo["memory"] = it.Spec.Memory.Guest.String()
		}
		instancetypeList = append(instancetypeList, itInfo)
	}

	jsonData, err := json.MarshalIndent(instancetypeList, "", "  ")
	if err != nil {
		return nil, err
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(jsonData),
			},
		},
	}, nil
}

func ClusterPreferencesList(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	parts := strings.Split(req.Params.URI, "/")
	if len(parts) < 3 || parts[2] != "cluster" {
		return nil, fmt.Errorf("invalid URI format, expected kubevirt://cluster/preferences")
	}

	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return nil, err
	}

	preferences, err := virtClient.VirtualMachineClusterPreference().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	preferenceList := make([]map[string]interface{}, 0, len(preferences.Items))
	for _, pref := range preferences.Items {
		prefInfo := map[string]interface{}{
			"name":    pref.Name,
			"created": pref.CreationTimestamp,
		}
		if pref.Spec.Machine != nil {
			prefInfo["hasMachinePreferences"] = true
		}
		preferenceList = append(preferenceList, prefInfo)
	}

	jsonData, err := json.MarshalIndent(preferenceList, "", "  ")
	if err != nil {
		return nil, err
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(jsonData),
			},
		},
	}, nil
}

func ClusterInstancetypeGet(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	parts := strings.Split(req.Params.URI, "/")
	if len(parts) < 5 || parts[2] != "cluster" {
		return nil, fmt.Errorf("invalid URI format, expected kubevirt://cluster/instancetype/{name}")
	}
	name := parts[4]
	if name == "" {
		return nil, fmt.Errorf("resource name may not be empty")
	}

	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return nil, err
	}

	instancetype, err := virtClient.VirtualMachineClusterInstancetype().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	jsonData, err := json.MarshalIndent(instancetype, "", "  ")
	if err != nil {
		return nil, err
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(jsonData),
			},
		},
	}, nil
}

func ClusterPreferenceGet(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	parts := strings.Split(req.Params.URI, "/")
	if len(parts) < 5 || parts[2] != "cluster" {
		return nil, fmt.Errorf("invalid URI format, expected kubevirt://cluster/preference/{name}")
	}
	name := parts[4]
	if name == "" {
		return nil, fmt.Errorf("resource name may not be empty")
	}

	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return nil, err
	}

	preference, err := virtClient.VirtualMachineClusterPreference().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	jsonData, err := json.MarshalIndent(preference, "", "  ")
	if err != nil {
		return nil, err
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(jsonData),
			},
		},
	}, nil
}