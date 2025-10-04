package vm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GetInstancetypeInput struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type GetInstancetypeOutput struct {
	Message string `json:"message"`
}

func GetInstancetype(ctx context.Context, req *mcp.CallToolRequest, input GetInstancetypeInput) (*mcp.CallToolResult, *GetInstancetypeOutput, error) {
	if input.Namespace == "" {
		return nil, nil, fmt.Errorf("namespace parameter is required")
	}
	if input.Name == "" {
		return nil, nil, fmt.Errorf("name parameter is required")
	}

	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return nil, nil, err
	}

	vm, err := virtClient.VirtualMachine(input.Namespace).Get(ctx, input.Name, metav1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}

	message := "no instance type referenced by virtual machine"
	if vm.Spec.Instancetype != nil {
		message = vm.Spec.Instancetype.Name
	}
	return nil, &GetInstancetypeOutput{Message: message}, nil
}

type GetStatusInput struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type GetStatusOutput struct {
	Result string `json:"result"`
}

func GetStatus(ctx context.Context, req *mcp.CallToolRequest, input GetStatusInput) (*mcp.CallToolResult, *GetStatusOutput, error) {
	if input.Namespace == "" {
		return nil, nil, fmt.Errorf("namespace parameter is required")
	}
	if input.Name == "" {
		return nil, nil, fmt.Errorf("name parameter is required")
	}

	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return nil, nil, err
	}

	vm, err := virtClient.VirtualMachine(input.Namespace).Get(ctx, input.Name, metav1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}

	statusInfo := map[string]interface{}{
		"name":               vm.Name,
		"namespace":          vm.Namespace,
		"status":             vm.Status.PrintableStatus,
		"ready":              vm.Status.Ready,
		"created":            vm.CreationTimestamp,
		"desiredGeneration":  vm.Status.DesiredGeneration,
		"observedGeneration": vm.Status.ObservedGeneration,
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

	resultJSON, err := json.MarshalIndent(statusInfo, "", "  ")
	if err != nil {
		return nil, nil, err
	}

	return nil, &GetStatusOutput{Result: string(resultJSON)}, nil
}

type GetConditionsInput struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type GetConditionsOutput struct {
	Result string `json:"result"`
}

func GetConditions(ctx context.Context, req *mcp.CallToolRequest, input GetConditionsInput) (*mcp.CallToolResult, *GetConditionsOutput, error) {
	if input.Namespace == "" {
		return nil, nil, fmt.Errorf("namespace parameter is required")
	}
	if input.Name == "" {
		return nil, nil, fmt.Errorf("name parameter is required")
	}

	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return nil, nil, err
	}

	vm, err := virtClient.VirtualMachine(input.Namespace).Get(ctx, input.Name, metav1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}

	conditionsInfo := map[string]interface{}{
		"name":       vm.Name,
		"namespace":  vm.Namespace,
		"conditions": []map[string]interface{}{},
	}

	if len(vm.Status.Conditions) > 0 {
		conditions := make([]map[string]interface{}, 0, len(vm.Status.Conditions))
		for _, cond := range vm.Status.Conditions {
			condition := map[string]interface{}{
				"type":               cond.Type,
				"status":             cond.Status,
				"lastTransitionTime": cond.LastTransitionTime,
			}
			if cond.Reason != "" {
				condition["reason"] = cond.Reason
			}
			if cond.Message != "" {
				condition["message"] = cond.Message
			}
			conditions = append(conditions, condition)
		}
		conditionsInfo["conditions"] = conditions
	}

	resultJSON, err := json.MarshalIndent(conditionsInfo, "", "  ")
	if err != nil {
		return nil, nil, err
	}

	return nil, &GetConditionsOutput{Result: string(resultJSON)}, nil
}

type GetPhaseInput struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type GetPhaseOutput struct {
	Result string `json:"result"`
}

func GetPhase(ctx context.Context, req *mcp.CallToolRequest, input GetPhaseInput) (*mcp.CallToolResult, *GetPhaseOutput, error) {
	if input.Namespace == "" {
		return nil, nil, fmt.Errorf("namespace parameter is required")
	}
	if input.Name == "" {
		return nil, nil, fmt.Errorf("name parameter is required")
	}

	virtClient, err := client.GetKubevirtClient()
	if err != nil {
		return nil, nil, err
	}

	vm, err := virtClient.VirtualMachine(input.Namespace).Get(ctx, input.Name, metav1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}

	phaseInfo := map[string]interface{}{
		"name":      vm.Name,
		"namespace": vm.Namespace,
		"status":    vm.Status.PrintableStatus,
		"ready":     vm.Status.Ready,
	}

	if vm.Spec.RunStrategy != nil {
		phaseInfo["runStrategy"] = string(*vm.Spec.RunStrategy)
	}

	resultJSON, err := json.MarshalIndent(phaseInfo, "", "  ")
	if err != nil {
		return nil, nil, err
	}

	return nil, &GetPhaseOutput{Result: string(resultJSON)}, nil
}