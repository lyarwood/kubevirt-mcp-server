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

		Context("when given valid arguments", func() {
			It("should accept valid namespace parameter", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "test-namespace",
				}

				// This will fail due to no KubeVirt cluster, but we're testing the argument parsing
				result, err := vm.List(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				// Should not contain argument parsing errors
				Expect(result.Content[0].(mcp.TextContent).Text).NotTo(ContainSubstring("parameter required"))
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

		Context("when given valid arguments", func() {
			It("should accept valid namespace and name parameters", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "test-namespace",
					"name":      "test-vm",
				}

				// This will fail due to no KubeVirt cluster, but we're testing the argument parsing
				result, err := vm.Start(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				// Should not contain argument parsing errors
				Expect(result.Content[0].(mcp.TextContent).Text).NotTo(ContainSubstring("parameter required"))
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

		Context("when given valid arguments", func() {
			It("should accept valid namespace and name parameters", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "test-namespace",
					"name":      "test-vm",
				}

				// This will fail due to no KubeVirt cluster, but we're testing the argument parsing
				result, err := vm.Stop(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				// Should not contain argument parsing errors
				Expect(result.Content[0].(mcp.TextContent).Text).NotTo(ContainSubstring("parameter required"))
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

		Context("when given valid arguments", func() {
			It("should accept valid namespace and name parameters", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "test-namespace",
					"name":      "test-vm",
				}

				// This will fail due to no KubeVirt cluster, but we're testing the argument parsing
				result, err := vm.Restart(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				// Should not contain argument parsing errors
				Expect(result.Content[0].(mcp.TextContent).Text).NotTo(ContainSubstring("parameter required"))
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

		Context("when given valid arguments", func() {
			It("should accept valid namespace and name parameters", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "test-namespace",
					"name":      "test-vm",
				}

				// This will fail due to no KubeVirt cluster, but we're testing the argument parsing
				result, err := vm.GetInstancetype(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				// Should not contain argument parsing errors
				Expect(result.Content[0].(mcp.TextContent).Text).NotTo(ContainSubstring("parameter required"))
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

		Context("when given valid arguments with container disk resolution", func() {
			It("should accept valid arguments and resolve container disk", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace":      "test-namespace",
					"name":           "test-vm",
					"container_disk": "fedora",
				}

				// This will fail due to no KubeVirt cluster, but we're testing the argument parsing and container disk resolution
				result, err := vm.Create(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				// Should not contain argument parsing errors
				Expect(result.Content[0].(mcp.TextContent).Text).NotTo(ContainSubstring("parameter required"))
			})
		})
	})

	Describe("Delete", func() {
		Context("when given invalid arguments", func() {
			It("should return an error for missing namespace", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"name": "test-vm",
				}

				result, err := vm.Delete(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("namespace parameter required"))
			})

			It("should return an error for missing name", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "test-ns",
				}

				result, err := vm.Delete(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("name parameter required"))
			})
		})

		Context("when given valid arguments", func() {
			It("should accept valid namespace and name parameters", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "test-namespace",
					"name":      "test-vm",
				}

				// This will fail due to no KubeVirt cluster, but we're testing the argument parsing
				result, err := vm.Delete(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				// Should not contain argument parsing errors
				Expect(result.Content[0].(mcp.TextContent).Text).NotTo(ContainSubstring("parameter required"))
			})
		})
	})

	Describe("Pause", func() {
		Context("when given invalid arguments", func() {
			It("should return an error for missing namespace", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"name": "test-vm",
				}

				result, err := vm.Pause(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("namespace parameter required"))
			})

			It("should return an error for missing name", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "test-ns",
				}

				result, err := vm.Pause(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("name parameter required"))
			})
		})

		Context("when given valid arguments", func() {
			It("should accept valid namespace and name parameters", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "test-namespace",
					"name":      "test-vm",
				}

				// This will fail due to no KubeVirt cluster, but we're testing the argument parsing
				result, err := vm.Pause(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				// Should not contain argument parsing errors
				Expect(result.Content[0].(mcp.TextContent).Text).NotTo(ContainSubstring("parameter required"))
			})
		})
	})

	Describe("Unpause", func() {
		Context("when given invalid arguments", func() {
			It("should return an error for missing namespace", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"name": "test-vm",
				}

				result, err := vm.Unpause(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("namespace parameter required"))
			})

			It("should return an error for missing name", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "test-ns",
				}

				result, err := vm.Unpause(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("name parameter required"))
			})
		})

		Context("when given valid arguments", func() {
			It("should accept valid namespace and name parameters", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "test-namespace",
					"name":      "test-vm",
				}

				// This will fail due to no KubeVirt cluster, but we're testing the argument parsing
				result, err := vm.Unpause(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				// Should not contain argument parsing errors
				Expect(result.Content[0].(mcp.TextContent).Text).NotTo(ContainSubstring("parameter required"))
			})
		})
	})

	Describe("GetStatus", func() {
		Context("when called with valid arguments", func() {
			It("should accept valid namespace and name parameters", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "default",
					"name":      "test-vm",
				}

				// This will fail due to no KubeVirt cluster, but we're testing the argument parsing
				result, err := vm.GetStatus(ctx, request)

				// We expect either no error (if mocked) or error at client creation stage
				if err != nil {
					Expect(result.IsError).To(BeTrue())
					// Should not contain argument parsing errors
					Expect(result.Content[0].(mcp.TextContent).Text).NotTo(ContainSubstring("parameter required"))
				} else {
					// If no error, the function succeeded in parsing arguments
					Expect(result.IsError).To(BeFalse())
				}
			})
		})

		Context("when called with invalid arguments", func() {
			It("should return an error for missing namespace", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"name": "test-vm",
				}

				result, err := vm.GetStatus(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("namespace parameter required"))
			})

			It("should return an error for missing name", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "default",
				}

				result, err := vm.GetStatus(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("name parameter required"))
			})
		})
	})

	Describe("GetConditions", func() {
		Context("when called with valid arguments", func() {
			It("should accept valid namespace and name parameters", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "default",
					"name":      "test-vm",
				}

				// This will fail due to no KubeVirt cluster, but we're testing the argument parsing
				result, err := vm.GetConditions(ctx, request)

				// We expect either no error (if mocked) or error at client creation stage
				if err != nil {
					Expect(result.IsError).To(BeTrue())
					// Should not contain argument parsing errors
					Expect(result.Content[0].(mcp.TextContent).Text).NotTo(ContainSubstring("parameter required"))
				} else {
					// If no error, the function succeeded in parsing arguments
					Expect(result.IsError).To(BeFalse())
				}
			})
		})

		Context("when called with invalid arguments", func() {
			It("should return an error for missing namespace", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"name": "test-vm",
				}

				result, err := vm.GetConditions(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("namespace parameter required"))
			})

			It("should return an error for missing name", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "default",
				}

				result, err := vm.GetConditions(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("name parameter required"))
			})
		})
	})

	Describe("GetPhase", func() {
		Context("when called with valid arguments", func() {
			It("should accept valid namespace and name parameters", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "default",
					"name":      "test-vm",
				}

				// This will fail due to no KubeVirt cluster, but we're testing the argument parsing
				result, err := vm.GetPhase(ctx, request)

				// We expect either no error (if mocked) or error at client creation stage
				if err != nil {
					Expect(result.IsError).To(BeTrue())
					// Should not contain argument parsing errors
					Expect(result.Content[0].(mcp.TextContent).Text).NotTo(ContainSubstring("parameter required"))
				} else {
					// If no error, the function succeeded in parsing arguments
					Expect(result.IsError).To(BeFalse())
				}
			})
		})

		Context("when called with invalid arguments", func() {
			It("should return an error for missing namespace", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"name": "test-vm",
				}

				result, err := vm.GetPhase(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("namespace parameter required"))
			})

			It("should return an error for missing name", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "default",
				}

				result, err := vm.GetPhase(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("name parameter required"))
			})
		})
	})

	Describe("Patch", func() {
		Context("when called with valid arguments", func() {
			It("should accept valid namespace, name, and patch parameters", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "default",
					"name":      "test-vm",
					"patch":     `{"metadata":{"labels":{"test":"label"}}}`,
				}

				// This will fail due to no KubeVirt cluster, but we're testing the argument parsing
				result, err := vm.Patch(ctx, request)

				// We expect either no error (if mocked) or error at client creation stage
				if err != nil {
					Expect(result.IsError).To(BeTrue())
					// Should not contain argument parsing errors
					Expect(result.Content[0].(mcp.TextContent).Text).NotTo(ContainSubstring("parameter required"))
					// Should not contain JSON validation errors for valid JSON
					Expect(result.Content[0].(mcp.TextContent).Text).NotTo(ContainSubstring("invalid JSON"))
				} else {
					// If no error, the function succeeded in parsing arguments
					Expect(result.IsError).To(BeFalse())
				}
			})
		})

		Context("when called with invalid arguments", func() {
			It("should return an error for missing namespace", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"name":  "test-vm",
					"patch": `{"metadata":{"labels":{"test":"label"}}}`,
				}

				result, err := vm.Patch(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("namespace parameter required"))
			})

			It("should return an error for missing name", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "default",
					"patch":     `{"metadata":{"labels":{"test":"label"}}}`,
				}

				result, err := vm.Patch(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("name parameter required"))
			})

			It("should return an error for missing patch", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "default",
					"name":      "test-vm",
				}

				result, err := vm.Patch(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("patch parameter required"))
			})

			It("should return an error for invalid JSON patch", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "default",
					"name":      "test-vm",
					"patch":     `{invalid json}`,
				}

				result, err := vm.Patch(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("invalid JSON in patch parameter"))
			})
		})
	})

	Describe("Disks", func() {
		Context("when given invalid arguments", func() {
			It("should return an error for missing namespace", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"name": "test-vm",
				}

				result, err := vm.Disks(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("namespace parameter required"))
			})

			It("should return an error for missing name", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "test-ns",
				}

				result, err := vm.Disks(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("name parameter required"))
			})
		})

		Context("when given valid arguments", func() {
			It("should accept valid namespace and name parameters", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "test-namespace",
					"name":      "test-vm",
				}

				// This will fail due to no KubeVirt cluster, but we're testing the argument parsing
				result, err := vm.Disks(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				// Should not contain argument parsing errors
				Expect(result.Content[0].(mcp.TextContent).Text).NotTo(ContainSubstring("parameter required"))
			})
		})
	})
})
