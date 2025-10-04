package vm_test

import (
	"context"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/tools/vm"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("VM", func() {
	var (
		ctx context.Context
		req *mcp.CallToolRequest
	)

	BeforeEach(func() {
		ctx = context.Background()
		req = &mcp.CallToolRequest{}
	})

	DescribeTable("testing vm tools with missing arguments",
		func(toolFunc interface{}, input interface{}, expectedErr string) {
			var err error
			switch f := toolFunc.(type) {
			case func(context.Context, *mcp.CallToolRequest, vm.ListInput) (*mcp.CallToolResult, *vm.ListOutput, error):
				_, _, err = f(ctx, req, input.(vm.ListInput))
			case func(context.Context, *mcp.CallToolRequest, vm.StartInput) (*mcp.CallToolResult, *vm.StartOutput, error):
				_, _, err = f(ctx, req, input.(vm.StartInput))
			case func(context.Context, *mcp.CallToolRequest, vm.StopInput) (*mcp.CallToolResult, *vm.StopOutput, error):
				_, _, err = f(ctx, req, input.(vm.StopInput))
			case func(context.Context, *mcp.CallToolRequest, vm.RestartInput) (*mcp.CallToolResult, *vm.RestartOutput, error):
				_, _, err = f(ctx, req, input.(vm.RestartInput))
			case func(context.Context, *mcp.CallToolRequest, vm.GetInstancetypeInput) (*mcp.CallToolResult, *vm.GetInstancetypeOutput, error):
				_, _, err = f(ctx, req, input.(vm.GetInstancetypeInput))
			case func(context.Context, *mcp.CallToolRequest, vm.CreateInput) (*mcp.CallToolResult, *vm.CreateOutput, error):
				_, _, err = f(ctx, req, input.(vm.CreateInput))
			case func(context.Context, *mcp.CallToolRequest, vm.DeleteInput) (*mcp.CallToolResult, *vm.DeleteOutput, error):
				_, _, err = f(ctx, req, input.(vm.DeleteInput))
			case func(context.Context, *mcp.CallToolRequest, vm.PauseInput) (*mcp.CallToolResult, *vm.PauseOutput, error):
				_, _, err = f(ctx, req, input.(vm.PauseInput))
			case func(context.Context, *mcp.CallToolRequest, vm.UnpauseInput) (*mcp.CallToolResult, *vm.UnpauseOutput, error):
				_, _, err = f(ctx, req, input.(vm.UnpauseInput))
			case func(context.Context, *mcp.CallToolRequest, vm.GetStatusInput) (*mcp.CallToolResult, *vm.GetStatusOutput, error):
				_, _, err = f(ctx, req, input.(vm.GetStatusInput))
			case func(context.Context, *mcp.CallToolRequest, vm.GetConditionsInput) (*mcp.CallToolResult, *vm.GetConditionsOutput, error):
				_, _, err = f(ctx, req, input.(vm.GetConditionsInput))
			case func(context.Context, *mcp.CallToolRequest, vm.GetPhaseInput) (*mcp.CallToolResult, *vm.GetPhaseOutput, error):
				_, _, err = f(ctx, req, input.(vm.GetPhaseInput))
			case func(context.Context, *mcp.CallToolRequest, vm.PatchInput) (*mcp.CallToolResult, *vm.PatchOutput, error):
				_, _, err = f(ctx, req, input.(vm.PatchInput))
			case func(context.Context, *mcp.CallToolRequest, vm.DisksInput) (*mcp.CallToolResult, *vm.DisksOutput, error):
				_, _, err = f(ctx, req, input.(vm.DisksInput))
			}
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(expectedErr))
		},
		Entry("List should require namespace", vm.List, vm.ListInput{}, "namespace parameter is required"),
		Entry("Start should require namespace", vm.Start, vm.StartInput{Name: "test"}, "namespace parameter is required"),
		Entry("Start should require name", vm.Start, vm.StartInput{Namespace: "test"}, "name parameter is required"),
		Entry("Stop should require namespace", vm.Stop, vm.StopInput{Name: "test"}, "namespace parameter is required"),
		Entry("Stop should require name", vm.Stop, vm.StopInput{Namespace: "test"}, "name parameter is required"),
		Entry("Restart should require namespace", vm.Restart, vm.RestartInput{Name: "test"}, "namespace parameter is required"),
		Entry("Restart should require name", vm.Restart, vm.RestartInput{Namespace: "test"}, "name parameter is required"),
		Entry("GetInstancetype should require namespace", vm.GetInstancetype, vm.GetInstancetypeInput{Name: "test"}, "namespace parameter is required"),
		Entry("GetInstancetype should require name", vm.GetInstancetype, vm.GetInstancetypeInput{Namespace: "test"}, "name parameter is required"),
		Entry("Create should require namespace", vm.Create, vm.CreateInput{Name: "test", ContainerDisk: "d"}, "namespace parameter is required"),
		Entry("Create should require name", vm.Create, vm.CreateInput{Namespace: "test", ContainerDisk: "d"}, "name parameter is required"),
		Entry("Create should require container_disk", vm.Create, vm.CreateInput{Namespace: "test", Name: "test"}, "container_disk parameter is required"),
		Entry("Delete should require namespace", vm.Delete, vm.DeleteInput{Name: "test"}, "namespace parameter is required"),
		Entry("Delete should require name", vm.Delete, vm.DeleteInput{Namespace: "test"}, "name parameter is required"),
		Entry("Pause should require namespace", vm.Pause, vm.PauseInput{Name: "test"}, "namespace parameter is required"),
		Entry("Pause should require name", vm.Pause, vm.PauseInput{Namespace: "test"}, "name parameter is required"),
		Entry("Unpause should require namespace", vm.Unpause, vm.UnpauseInput{Name: "test"}, "namespace parameter is required"),
		Entry("Unpause should require name", vm.Unpause, vm.UnpauseInput{Namespace: "test"}, "name parameter is required"),
		Entry("GetStatus should require namespace", vm.GetStatus, vm.GetStatusInput{Name: "test"}, "namespace parameter is required"),
		Entry("GetStatus should require name", vm.GetStatus, vm.GetStatusInput{Namespace: "test"}, "name parameter is required"),
		Entry("GetConditions should require namespace", vm.GetConditions, vm.GetConditionsInput{Name: "test"}, "namespace parameter is required"),
		Entry("GetConditions should require name", vm.GetConditions, vm.GetConditionsInput{Namespace: "test"}, "name parameter is required"),
		Entry("GetPhase should require namespace", vm.GetPhase, vm.GetPhaseInput{Name: "test"}, "namespace parameter is required"),
		Entry("GetPhase should require name", vm.GetPhase, vm.GetPhaseInput{Namespace: "test"}, "name parameter is required"),
		Entry("Patch should require namespace", vm.Patch, vm.PatchInput{Name: "test", Patch: "{}"}, "namespace parameter is required"),
		Entry("Patch should require name", vm.Patch, vm.PatchInput{Namespace: "test", Patch: "{}"}, "name parameter is required"),
		Entry("Patch should require patch", vm.Patch, vm.PatchInput{Namespace: "test", Name: "test"}, "patch parameter is required"),
		Entry("Patch should require valid json patch", vm.Patch, vm.PatchInput{Namespace: "test", Name: "test", Patch: "bad"}, "invalid JSON in patch parameter"),
		Entry("Disks should require namespace", vm.Disks, vm.DisksInput{Name: "test"}, "namespace parameter is required"),
		Entry("Disks should require name", vm.Disks, vm.DisksInput{Namespace: "test"}, "name parameter is required"),
	)
})