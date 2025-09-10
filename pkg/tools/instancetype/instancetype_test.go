package instancetype_test

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/tools/instancetype"
)

var _ = Describe("Instancetype", func() {
	var (
		ctx context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
	})

	Describe("List", func() {
		Context("when called with valid request", func() {
			It("should accept empty arguments", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{}

				// This will fail due to no KubeVirt cluster, but we're testing the argument parsing
				result, err := instancetype.List(ctx, request)

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

			It("should accept nil arguments", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = nil

				// This will fail due to no KubeVirt cluster, but we're testing the argument parsing
				result, err := instancetype.List(ctx, request)

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

			It("should accept extra arguments gracefully", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"extra_param": "should_be_ignored",
				}

				// This will fail due to no KubeVirt cluster, but we're testing the argument parsing
				result, err := instancetype.List(ctx, request)

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

		Context("when client creation fails", func() {
			It("should handle client creation errors properly", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{}

				// This will fail due to no KubeVirt cluster - testing error handling path
				result, err := instancetype.List(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content).To(HaveLen(1))
				Expect(result.Content[0]).To(BeAssignableToTypeOf(mcp.TextContent{}))
			})
		})
	})
})