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
		})
	})
})