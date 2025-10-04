package vm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/client"
	"github.com/mark3labs/mcp-go/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

func GetStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	// Create a comprehensive status response
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

	// Add state change requests if available
	if len(vm.Status.StateChangeRequests) > 0 {
		requests := make([]map[string]interface{}, 0, len(vm.Status.StateChangeRequests))
		for _, req := range vm.Status.StateChangeRequests {
			request := map[string]interface{}{
				"action": req.Action,
			}
			if req.UID != nil {
				request["uid"] = *req.UID
			}
			requests = append(requests, request)
		}
		statusInfo["stateChangeRequests"] = requests
	}

	resultJSON, err := json.MarshalIndent(statusInfo, "", "  ")
	if err != nil {
		return newToolResultErr(err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(resultJSON),
			},
		},
	}, nil
}

func GetConditions(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	// Create conditions response
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
		return newToolResultErr(err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(resultJSON),
			},
		},
	}, nil
}

func GetPhase(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	// Simple phase response
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
		return newToolResultErr(err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(resultJSON),
			},
		},
	}, nil
}
