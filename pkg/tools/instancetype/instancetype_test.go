package instancetype_test

import (
	"context"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/tools/instancetype"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Instancetype", func() {
	var (
		ctx context.Context
		req *mcp.CallToolRequest
	)

	BeforeEach(func() {
		ctx = context.Background()
		req = &mcp.CallToolRequest{}
	})

	Describe("List", func() {
		Context("when called", func() {
			It("should proceed past argument validation", func() {
				input := instancetype.ListInput{}
				_, _, err := instancetype.List(ctx, req, input)
				// Expect an error from client creation, not from argument parsing.
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).NotTo(ContainSubstring("parameter"))
			})
		})
	})

	Describe("Get", func() {
		Context("when called with valid arguments", func() {
			It("should proceed past argument validation", func() {
				input := instancetype.GetInput{Name: "m1.large"}
				_, _, err := instancetype.Get(ctx, req, input)
				// Expect an error from client creation, not from argument parsing.
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).NotTo(ContainSubstring("name parameter is required"))
			})
		})

		Context("when called with invalid arguments", func() {
			It("should reject missing name parameter", func() {
				input := instancetype.GetInput{}
				_, _, err := instancetype.Get(ctx, req, input)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("name parameter is required"))
			})
		})
	})
})