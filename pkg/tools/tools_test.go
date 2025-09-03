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
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("namespace parameter required"))
			})

			It("should return an error for non-string namespace", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": 123,
				}

				result, err := tools.VmsList(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("namespace parameter required"))
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
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("namespace parameter required"))
			})

			It("should return an error for missing name", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "test-ns",
				}

				result, err := tools.VmStart(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("name parameter required"))
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
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("namespace parameter required"))
			})

			It("should return an error for missing name", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "test-ns",
				}

				result, err := tools.VmStop(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("name parameter required"))
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
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("namespace parameter required"))
			})

			It("should return an error for missing name", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "test-ns",
				}

				result, err := tools.VmRestart(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("name parameter required"))
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
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("namespace parameter required"))
			})

			It("should return an error for missing name", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace": "test-ns",
				}

				result, err := tools.VmGetInstancetype(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("name parameter required"))
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

				// We expect either no error (if mocked) or error at client creation stage
				if err != nil {
					Expect(result.IsError).To(BeTrue())
					// Should not contain argument parsing errors
					Expect(result.Content[0].(mcp.TextContent).Text).NotTo(ContainSubstring("unable to decode"))
				} else {
					// If no error, the function succeeded in parsing arguments
					Expect(result.IsError).To(BeFalse())
				}
			})
		})
	})

	Describe("VmCreate", func() {
		Context("when given invalid arguments", func() {
			It("should return an error for missing namespace", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"name":           "test-vm",
					"container_disk": "quay.io/kubevirt/cirros-container-disk-demo",
				}

				result, err := tools.VmCreate(ctx, request)

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

				result, err := tools.VmCreate(ctx, request)

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

				result, err := tools.VmCreate(ctx, request)

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
				result, err := tools.VmCreate(ctx, request)

				// We expect it to fail at the client creation stage, not argument parsing
				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				// Should not contain argument parsing errors
				Expect(result.Content[0].(mcp.TextContent).Text).NotTo(ContainSubstring("unable to decode"))
			})

			It("should accept optional instancetype and preference arguments", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"namespace":      "test-ns",
					"name":           "test-vm",
					"container_disk": "quay.io/kubevirt/cirros-container-disk-demo",
					"instancetype":   "small",
					"preference":     "fedora",
				}

				// This will fail due to no KubeVirt cluster, but we're testing the argument parsing
				result, err := tools.VmCreate(ctx, request)

				// We expect it to fail at the client creation stage, not argument parsing
				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				// Should not contain argument parsing errors
				Expect(result.Content[0].(mcp.TextContent).Text).NotTo(ContainSubstring("unable to decode"))
			})
		})
	})

	Describe("ResolveContainerDisk", func() {
		Context("when given OS names", func() {
			It("should resolve fedora to containerdisks fedora image", func() {
				result := tools.ResolveContainerDisk("fedora")
				Expect(result).To(Equal("quay.io/containerdisks/fedora:latest"))
			})

			It("should resolve ubuntu to containerdisks ubuntu image", func() {
				result := tools.ResolveContainerDisk("ubuntu")
				Expect(result).To(Equal("quay.io/containerdisks/ubuntu:latest"))
			})

			It("should resolve centos to containerdisks centos image", func() {
				result := tools.ResolveContainerDisk("centos")
				Expect(result).To(Equal("quay.io/containerdisks/centos:latest"))
			})

			It("should resolve cirros to kubevirt demo image", func() {
				result := tools.ResolveContainerDisk("cirros")
				Expect(result).To(Equal("quay.io/kubevirt/cirros-container-disk-demo"))
			})

			It("should handle case insensitive input", func() {
				result := tools.ResolveContainerDisk("FEDORA")
				Expect(result).To(Equal("quay.io/containerdisks/fedora:latest"))
			})

			It("should handle input with extra whitespace", func() {
				result := tools.ResolveContainerDisk("  ubuntu  ")
				Expect(result).To(Equal("quay.io/containerdisks/ubuntu:latest"))
			})
		})

		Context("when given container image names", func() {
			It("should return full container image URLs as-is", func() {
				input := "quay.io/containerdisks/fedora:38"
				result := tools.ResolveContainerDisk(input)
				Expect(result).To(Equal(input))
			})

			It("should return images with tags as-is", func() {
				input := "my-registry/my-image:v1.0"
				result := tools.ResolveContainerDisk(input)
				Expect(result).To(Equal(input))
			})

			It("should return images with slashes as-is", func() {
				input := "docker.io/library/ubuntu"
				result := tools.ResolveContainerDisk(input)
				Expect(result).To(Equal(input))
			})
		})

		Context("when given unknown OS names", func() {
			It("should construct containerdisks URL for unknown OS", func() {
				result := tools.ResolveContainerDisk("myos")
				Expect(result).To(Equal("quay.io/containerdisks/myos:latest"))
			})

			It("should handle unknown OS with case normalization", func() {
				result := tools.ResolveContainerDisk("MyCustomOS")
				Expect(result).To(Equal("quay.io/containerdisks/mycustomos:latest"))
			})
		})
	})
})
