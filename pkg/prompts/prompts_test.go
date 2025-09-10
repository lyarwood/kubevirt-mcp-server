package prompts_test

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/prompts"
)

var _ = Describe("Prompts", func() {
	var (
		ctx context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
	})

	Describe("DescribeVM", func() {
		Context("when called with valid arguments", func() {
			It("should accept valid namespace and name parameters", func() {
				request := mcp.GetPromptRequest{}
				request.Params.Arguments = map[string]string{
					"namespace": "default",
					"name":      "test-vm",
				}

				result, err := prompts.DescribeVM(ctx, request)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.Description).To(ContainSubstring("Comprehensive description of virtual machine test-vm in namespace default"))
				Expect(result.Messages).To(HaveLen(1))
				Expect(result.Messages[0].Role).To(Equal(mcp.RoleUser))
				Expect(result.Messages[0].Content.(mcp.TextContent).Text).To(ContainSubstring("test-vm"))
				Expect(result.Messages[0].Content.(mcp.TextContent).Text).To(ContainSubstring("default"))
			})
		})

		Context("when called with invalid arguments", func() {
			It("should return an error for missing namespace", func() {
				request := mcp.GetPromptRequest{}
				request.Params.Arguments = map[string]string{
					"name": "test-vm",
				}

				result, err := prompts.DescribeVM(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(err.Error()).To(ContainSubstring("namespace parameter is required"))
			})

			It("should return an error for missing name", func() {
				request := mcp.GetPromptRequest{}
				request.Params.Arguments = map[string]string{
					"namespace": "default",
				}

				result, err := prompts.DescribeVM(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(err.Error()).To(ContainSubstring("name parameter is required"))
			})
		})
	})

	Describe("TroubleshootVM", func() {
		Context("when called with valid arguments", func() {
			It("should accept valid namespace and name parameters", func() {
				request := mcp.GetPromptRequest{}
				request.Params.Arguments = map[string]string{
					"namespace": "default",
					"name":      "test-vm",
				}

				result, err := prompts.TroubleshootVM(ctx, request)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.Description).To(ContainSubstring("Comprehensive troubleshooting analysis for virtual machine test-vm in namespace default"))
				Expect(result.Messages).To(HaveLen(1))
				Expect(result.Messages[0].Role).To(Equal(mcp.RoleUser))
				Expect(result.Messages[0].Content.(mcp.TextContent).Text).To(ContainSubstring("test-vm"))
				Expect(result.Messages[0].Content.(mcp.TextContent).Text).To(ContainSubstring("default"))
			})

			It("should include issue description when provided", func() {
				request := mcp.GetPromptRequest{}
				request.Params.Arguments = map[string]string{
					"namespace":         "default",
					"name":              "test-vm",
					"issue_description": "VM won't start",
				}

				result, err := prompts.TroubleshootVM(ctx, request)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.Messages[0].Content.(mcp.TextContent).Text).To(ContainSubstring("VM won't start"))
			})
		})

		Context("when called with invalid arguments", func() {
			It("should return an error for missing namespace", func() {
				request := mcp.GetPromptRequest{}
				request.Params.Arguments = map[string]string{
					"name": "test-vm",
				}

				result, err := prompts.TroubleshootVM(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(err.Error()).To(ContainSubstring("namespace parameter is required"))
			})

			It("should return an error for missing name", func() {
				request := mcp.GetPromptRequest{}
				request.Params.Arguments = map[string]string{
					"namespace": "default",
				}

				result, err := prompts.TroubleshootVM(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(err.Error()).To(ContainSubstring("name parameter is required"))
			})
		})
	})

	Describe("HealthCheckVM", func() {
		Context("when called with valid arguments", func() {
			It("should accept valid namespace and name parameters", func() {
				request := mcp.GetPromptRequest{}
				request.Params.Arguments = map[string]string{
					"namespace": "default",
					"name":      "test-vm",
				}

				result, err := prompts.HealthCheckVM(ctx, request)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.Description).To(ContainSubstring("Quick health assessment of virtual machine test-vm in namespace default"))
				Expect(result.Messages).To(HaveLen(1))
				Expect(result.Messages[0].Role).To(Equal(mcp.RoleUser))
				Expect(result.Messages[0].Content.(mcp.TextContent).Text).To(ContainSubstring("test-vm"))
				Expect(result.Messages[0].Content.(mcp.TextContent).Text).To(ContainSubstring("default"))
				Expect(result.Messages[0].Content.(mcp.TextContent).Text).To(ContainSubstring("Quick Health Check"))
			})
		})

		Context("when called with invalid arguments", func() {
			It("should return an error for missing namespace", func() {
				request := mcp.GetPromptRequest{}
				request.Params.Arguments = map[string]string{
					"name": "test-vm",
				}

				result, err := prompts.HealthCheckVM(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(err.Error()).To(ContainSubstring("namespace parameter is required"))
			})

			It("should return an error for missing name", func() {
				request := mcp.GetPromptRequest{}
				request.Params.Arguments = map[string]string{
					"namespace": "default",
				}

				result, err := prompts.HealthCheckVM(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(err.Error()).To(ContainSubstring("name parameter is required"))
			})
		})
	})
})