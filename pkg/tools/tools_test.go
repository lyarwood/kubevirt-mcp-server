package tools_test

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/tools"
)

var _ = Describe("Tools", func() {
	var (
		ctx context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
	})

	Describe("VmsList", func() {
		Context("when given invalid arguments", func() {
			It("should return an error for missing namespace", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{}

				result, err := tools.VmsList(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("unable to decode namespace string"))
			})

			It("should return an error for non-string namespace", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": 123,
				}

				result, err := tools.VmsList(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("unable to decode namespace string"))
			})
		})
	})

	Describe("VmStart", func() {
		Context("when given invalid arguments", func() {
			It("should return an error for missing namespace", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"name": "test-vm",
				}

				result, err := tools.VmStart(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("unable to decode namespace string"))
			})

			It("should return an error for missing name", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "test-ns",
				}

				result, err := tools.VmStart(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("unable to decode name string"))
			})
		})
	})

	Describe("VmStop", func() {
		Context("when given invalid arguments", func() {
			It("should return an error for missing namespace", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"name": "test-vm",
				}

				result, err := tools.VmStop(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("unable to decode namespace string"))
			})

			It("should return an error for missing name", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "test-ns",
				}

				result, err := tools.VmStop(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("unable to decode name string"))
			})
		})
	})

	Describe("VmRestart", func() {
		Context("when given invalid arguments", func() {
			It("should return an error for missing namespace", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"name": "test-vm",
				}

				result, err := tools.VmRestart(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("unable to decode namespace string"))
			})

			It("should return an error for missing name", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "test-ns",
				}

				result, err := tools.VmRestart(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("unable to decode name string"))
			})
		})
	})

	Describe("VmGetInstancetype", func() {
		Context("when given invalid arguments", func() {
			It("should return an error for missing namespace", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"name": "test-vm",
				}

				result, err := tools.VmGetInstancetype(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("unable to decode namespace string"))
			})

			It("should return an error for missing name", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "test-ns",
				}

				result, err := tools.VmGetInstancetype(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("unable to decode name string"))
			})
		})
	})

	Describe("InstancetypesList", func() {
		Context("when called with valid request", func() {
			It("should accept empty arguments", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{}

				// This will fail due to no KubeVirt cluster, but we're testing the argument parsing
				result, err := tools.InstancetypesList(ctx, request)

				// We expect it to fail at the client creation stage, not argument parsing
				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				// Should not contain argument parsing errors
				Expect(result.Content[0].(mcp.TextContent).Text).NotTo(ContainSubstring("unable to decode"))
			})
		})
	})
})
