package preference_test

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/tools/preference"
)

var _ = Describe("Preference", func() {
	var (
		ctx context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
	})

	Describe("Get", func() {
		Context("when called with valid arguments", func() {
			It("should accept valid name parameter", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"name": "fedora",
				}

				// This will fail due to no KubeVirt cluster, but we're testing the argument parsing
				result, err := preference.Get(ctx, request)

				// We expect either no error (if mocked) or error at client creation stage
				if err != nil {
					Expect(result.IsError).To(BeTrue())
					// Should not contain argument parsing errors for valid name
					Expect(result.Content[0].(mcp.TextContent).Text).NotTo(ContainSubstring("name parameter is required"))
				} else {
					// If no error, the function succeeded in parsing arguments
					Expect(result.IsError).To(BeFalse())
				}
			})
		})

		Context("when called with invalid arguments", func() {
			It("should reject missing name parameter", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{}

				result, err := preference.Get(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("name parameter is required"))
			})

			It("should reject empty name parameter", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"name": "",
				}

				result, err := preference.Get(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("resource name may not be empty"))
			})

			It("should reject non-string name parameter", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"name": 123,
				}

				result, err := preference.Get(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result.IsError).To(BeTrue())
				Expect(result.Content[0].(mcp.TextContent).Text).To(ContainSubstring("name parameter is required"))
			})
		})

		Context("when client creation fails", func() {
			It("should handle client creation errors properly", func() {
				request := mcp.CallToolRequest{}
				request.Params.Arguments = map[string]interface{}{
					"name": "fedora",
				}

				// This will fail due to no KubeVirt cluster - testing error handling path
				result, err := preference.Get(ctx, request)

				// We expect either no error (if mocked) or error at client creation stage
				if err != nil {
					Expect(result.IsError).To(BeTrue())
					Expect(result.Content).To(HaveLen(1))
					Expect(result.Content[0]).To(BeAssignableToTypeOf(mcp.TextContent{}))
				} else {
					// If no error, the function succeeded
					Expect(result.IsError).To(BeFalse())
				}
			})
		})
	})
})
