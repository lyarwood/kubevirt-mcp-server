package vm_test

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/tools/vm"
)

var _ = Describe("VM", func() {
	var (
		ctx context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
	})

	Describe("List", func() {
		Context("when given invalid arguments", func() {
			It("should return an error for missing namespace", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{}

				result, err := vm.List(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("namespace parameter required"))
			})

			It("should return an error for non-string namespace", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": 123,
				}

				result, err := vm.List(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("namespace parameter required"))
			})
		})
	})

	Describe("Start", func() {
		Context("when given invalid arguments", func() {
			It("should return an error for missing namespace", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"name": "test-vm",
				}

				result, err := vm.Start(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("namespace parameter required"))
			})

			It("should return an error for missing name", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "test-ns",
				}

				result, err := vm.Start(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("name parameter required"))
			})
		})
	})

	Describe("Stop", func() {
		Context("when given invalid arguments", func() {
			It("should return an error for missing namespace", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"name": "test-vm",
				}

				result, err := vm.Stop(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("namespace parameter required"))
			})

			It("should return an error for missing name", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "test-ns",
				}

				result, err := vm.Stop(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("name parameter required"))
			})
		})
	})

	Describe("Restart", func() {
		Context("when given invalid arguments", func() {
			It("should return an error for missing namespace", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"name": "test-vm",
				}

				result, err := vm.Restart(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("namespace parameter required"))
			})

			It("should return an error for missing name", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "test-ns",
				}

				result, err := vm.Restart(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("name parameter required"))
			})
		})
	})

	Describe("GetInstancetype", func() {
		Context("when given invalid arguments", func() {
			It("should return an error for missing namespace", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"name": "test-vm",
				}

				result, err := vm.GetInstancetype(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("namespace parameter required"))
			})

			It("should return an error for missing name", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "test-ns",
				}

				result, err := vm.GetInstancetype(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("name parameter required"))
			})
		})
	})

	Describe("Create", func() {
		Context("when given invalid arguments", func() {
			It("should return an error for missing namespace", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"name":           "test-vm",
					"container_disk": "quay.io/kubevirt/cirros-container-disk-demo",
				}

				result, err := vm.Create(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("namespace parameter required"))
			})

			It("should return an error for missing name", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace":      "test-ns",
					"container_disk": "quay.io/kubevirt/cirros-container-disk-demo",
				}

				result, err := vm.Create(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("name parameter required"))
			})

			It("should return an error for missing container_disk", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "test-ns",
					"name":      "test-vm",
				}

				result, err := vm.Create(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("container_disk parameter required"))
			})
		})

		Context("when given valid arguments", func() {
			It("should accept all required arguments without instancetype and preference", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace":      "test-ns",
					"name":           "test-vm",
					"container_disk": "quay.io/kubevirt/cirros-container-disk-demo",
				}

				// This will fail due to no KubeVirt cluster, but we're testing the argument parsing
				result, err := vm.Create(ctx, request)

				// We expect error at client creation stage, not argument parsing
				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				// Should not contain argument parsing errors
				Expect(result.Content[0].(mcp.TextContent).Text).NotTo(ContainSubstring("parameter required"))
			})

			It("should accept optional instancetype and preference arguments", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace":      "test-ns",
					"name":           "test-vm",
					"container_disk": "quay.io/kubevirt/cirros-container-disk-demo",
					"instancetype":   "u1.medium",
					"preference":     "fedora",
				}

				// This will fail due to no KubeVirt cluster, but we're testing the argument parsing
				result, err := vm.Create(ctx, request)

				// We expect error at client creation stage, not argument parsing
				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				// Should not contain argument parsing errors
				Expect(result.Content[0].(mcp.TextContent).Text).NotTo(ContainSubstring("parameter required"))
			})
		})
	})
})