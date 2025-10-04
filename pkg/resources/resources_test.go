package resources_test

import (
	"context"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/resources"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Resources", func() {
	var (
		ctx context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
	})

	type resourceFunc func(context.Context, *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error)

	DescribeTable("invalid URI handling",
		func(handler resourceFunc, uri string, expectedError string) {
			request := &mcp.ReadResourceRequest{
				Params: &mcp.ReadResourceParams{
					URI: uri,
				},
			}
			result, err := handler(ctx, request)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(expectedError))
			Expect(result).To(BeNil())
		},
		Entry("VmsList should fail with malformed URI", resources.VmsList, "invalid-uri", "invalid URI format"),
		Entry("VmsList should fail with missing namespace", resources.VmsList, "kubevirt:////vms", "resource namespace may not be empty"),
		Entry("VmGet should fail with malformed URI", resources.VmGet, "invalid-uri", "invalid URI format"),
		Entry("VmGet should fail with missing name", resources.VmGet, "kubevirt://test-ns/vm/", "resource name may not be empty"),
		Entry("VmisList should fail with malformed URI", resources.VmisList, "invalid-uri", "invalid URI format"),
		Entry("VmisList should fail with missing namespace", resources.VmisList, "kubevirt:////vmis", "resource namespace may not be empty"),
		Entry("VmiGet should fail with malformed URI", resources.VmiGet, "invalid-uri", "invalid URI format"),
		Entry("VmiGet should fail with missing name", resources.VmiGet, "kubevirt://test-ns/vmi/", "resource name may not be empty"),
		Entry("DataVolumesList should fail with malformed URI", resources.DataVolumesList, "invalid-uri", "invalid URI format"),
		Entry("DataVolumesList should fail with missing namespace", resources.DataVolumesList, "kubevirt:////datavolumes", "resource namespace may not be empty"),
		Entry("DataVolumeGet should fail with malformed URI", resources.DataVolumeGet, "invalid-uri", "invalid URI format"),
		Entry("DataVolumeGet should fail with missing name", resources.DataVolumeGet, "kubevirt://test-ns/datavolume/", "resource name may not be empty"),
		Entry("VmGetStatus should fail with malformed URI", resources.VmGetStatus, "invalid-uri", "invalid URI format"),
		Entry("VmGetStatus should fail with missing name", resources.VmGetStatus, "kubevirt://test-ns/vm//status", "resource name may not be empty"),
	)

	DescribeTable("valid URI parsing",
		func(handler resourceFunc, uri string) {
			request := &mcp.ReadResourceRequest{
				Params: &mcp.ReadResourceParams{
					URI: uri,
				},
			}
			// This will fail due to no KubeVirt cluster, but we're testing URI parsing
			_, err := handler(ctx, request)
			Expect(err).To(HaveOccurred())
			// Should not contain URI parsing errors
			Expect(err.Error()).NotTo(ContainSubstring("invalid URI format"))
			Expect(err.Error()).NotTo(ContainSubstring("may not be empty"))
		},
		Entry("VmsList should parse namespace", resources.VmsList, "kubevirt://test-namespace/vms"),
		Entry("VmGet should parse namespace and name", resources.VmGet, "kubevirt://test-namespace/vm/test-vm"),
		Entry("VmisList should parse namespace", resources.VmisList, "kubevirt://test-namespace/vmis"),
		Entry("VmiGet should parse namespace and name", resources.VmiGet, "kubevirt://test-namespace/vmi/test-vmi"),
		Entry("DataVolumesList should parse namespace", resources.DataVolumesList, "kubevirt://test-namespace/datavolumes"),
		Entry("DataVolumeGet should parse namespace and name", resources.DataVolumeGet, "kubevirt://test-namespace/datavolume/test-dv"),
		Entry("VmGetStatus should parse namespace and name", resources.VmGetStatus, "kubevirt://test-namespace/vm/test-vm/status"),
	)
})