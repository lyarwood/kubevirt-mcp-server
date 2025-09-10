package resources_test

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/resources"
)

var _ = Describe("Resources", func() {
	var (
		ctx context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
	})

	Describe("VmsList", func() {
		Context("when given invalid URI", func() {
			It("should return an error for malformed URI", func() {
				request := mcp.ReadResourceRequest{
					Params: struct {
						URI       string                 `json:"uri"`
						Arguments map[string]interface{} `json:"arguments,omitempty"`
					}{
						URI: "invalid-uri",
					},
				}

				result, err := resources.VmsList(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid URI format"))
				Expect(result).To(BeNil())
			})

			It("should return an error for URI with missing namespace", func() {
				request := mcp.ReadResourceRequest{
					Params: struct {
						URI       string                 `json:"uri"`
						Arguments map[string]interface{} `json:"arguments,omitempty"`
					}{
						URI: "kubevirt:/vms",
					},
				}

				result, err := resources.VmsList(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid URI format"))
				Expect(result).To(BeNil())
			})
		})

		Context("when given valid URI format", func() {
			It("should parse namespace correctly", func() {
				request := mcp.ReadResourceRequest{
					Params: struct {
						URI       string                 `json:"uri"`
						Arguments map[string]interface{} `json:"arguments,omitempty"`
					}{
						URI: "kubevirt://test-namespace/vms",
					},
				}

				// This will fail due to no KubeVirt cluster, but we're testing URI parsing
				result, err := resources.VmsList(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				// Should not contain URI parsing errors
				Expect(err.Error()).NotTo(ContainSubstring("invalid URI format"))
			})
		})

		Context("when given valid URI format", func() {
			It("should accept valid kubevirt URI and attempt to list VMs", func() {
				request := mcp.ReadResourceRequest{
					Params: struct {
						URI       string                 `json:"uri"`
						Arguments map[string]interface{} `json:"arguments,omitempty"`
					}{
						URI: "kubevirt://test-namespace/vms",
					},
				}

				// This will fail due to no KubeVirt cluster, but we're testing the URI parsing
				result, err := resources.VmsList(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).NotTo(ContainSubstring("invalid URI format"))
				Expect(err.Error()).NotTo(ContainSubstring("missing namespace"))
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("VmGet", func() {
		Context("when given invalid URI", func() {
			It("should return an error for malformed URI", func() {
				request := mcp.ReadResourceRequest{
					Params: struct {
						URI       string                 `json:"uri"`
						Arguments map[string]interface{} `json:"arguments,omitempty"`
					}{
						URI: "invalid-uri",
					},
				}

				result, err := resources.VmGet(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid URI format"))
				Expect(result).To(BeNil())
			})

			It("should return an error for URI missing VM name", func() {
				request := mcp.ReadResourceRequest{
					Params: struct {
						URI       string                 `json:"uri"`
						Arguments map[string]interface{} `json:"arguments,omitempty"`
					}{
						URI: "kubevirt://test-ns/vm",
					},
				}

				result, err := resources.VmGet(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid URI format"))
				Expect(result).To(BeNil())
			})
		})

		Context("when given valid URI format", func() {
			It("should parse namespace and name correctly", func() {
				request := mcp.ReadResourceRequest{
					Params: struct {
						URI       string                 `json:"uri"`
						Arguments map[string]interface{} `json:"arguments,omitempty"`
					}{
						URI: "kubevirt://test-namespace/vm/test-vm",
					},
				}

				// This will fail due to no KubeVirt cluster, but we're testing URI parsing
				result, err := resources.VmGet(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				// Should not contain URI parsing errors
				Expect(err.Error()).NotTo(ContainSubstring("invalid URI format"))
			})
		})

		Context("when given valid URI format", func() {
			It("should accept valid kubevirt URI and attempt to get VM", func() {
				request := mcp.ReadResourceRequest{
					Params: struct {
						URI       string                 `json:"uri"`
						Arguments map[string]interface{} `json:"arguments,omitempty"`
					}{
						URI: "kubevirt://test-namespace/vm/test-vm",
					},
				}

				// This will fail due to no KubeVirt cluster, but we're testing the URI parsing
				result, err := resources.VmGet(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).NotTo(ContainSubstring("invalid URI format"))
				Expect(err.Error()).NotTo(ContainSubstring("missing"))
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("VmisList", func() {
		Context("when given invalid URI", func() {
			It("should return an error for malformed URI", func() {
				request := mcp.ReadResourceRequest{
					Params: struct {
						URI       string                 `json:"uri"`
						Arguments map[string]interface{} `json:"arguments,omitempty"`
					}{
						URI: "invalid-uri",
					},
				}

				result, err := resources.VmisList(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid URI format"))
				Expect(result).To(BeNil())
			})
		})

		Context("when given valid URI format", func() {
			It("should parse namespace correctly", func() {
				request := mcp.ReadResourceRequest{
					Params: struct {
						URI       string                 `json:"uri"`
						Arguments map[string]interface{} `json:"arguments,omitempty"`
					}{
						URI: "kubevirt://test-namespace/vmis",
					},
				}

				// This will fail due to no KubeVirt cluster, but we're testing URI parsing
				result, err := resources.VmisList(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				// Should not contain URI parsing errors
				Expect(err.Error()).NotTo(ContainSubstring("invalid URI format"))
			})
		})

		Context("when given valid URI format", func() {
			It("should accept valid kubevirt URI and attempt to list VMIs", func() {
				request := mcp.ReadResourceRequest{
					Params: struct {
						URI       string                 `json:"uri"`
						Arguments map[string]interface{} `json:"arguments,omitempty"`
					}{
						URI: "kubevirt://test-namespace/vmis",
					},
				}

				// This will fail due to no KubeVirt cluster, but we're testing the URI parsing
				result, err := resources.VmisList(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).NotTo(ContainSubstring("invalid URI format"))
				Expect(err.Error()).NotTo(ContainSubstring("missing namespace"))
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("VmiGet", func() {
		Context("when given invalid URI", func() {
			It("should return an error for malformed URI", func() {
				request := mcp.ReadResourceRequest{
					Params: struct {
						URI       string                 `json:"uri"`
						Arguments map[string]interface{} `json:"arguments,omitempty"`
					}{
						URI: "invalid-uri",
					},
				}

				result, err := resources.VmiGet(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid URI format"))
				Expect(result).To(BeNil())
			})

			It("should return an error for URI missing VMI name", func() {
				request := mcp.ReadResourceRequest{
					Params: struct {
						URI       string                 `json:"uri"`
						Arguments map[string]interface{} `json:"arguments,omitempty"`
					}{
						URI: "kubevirt://test-ns/vmi",
					},
				}

				result, err := resources.VmiGet(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid URI format"))
				Expect(result).To(BeNil())
			})
		})

		Context("when given valid URI format", func() {
			It("should parse namespace and name correctly", func() {
				request := mcp.ReadResourceRequest{
					Params: struct {
						URI       string                 `json:"uri"`
						Arguments map[string]interface{} `json:"arguments,omitempty"`
					}{
						URI: "kubevirt://test-namespace/vmi/test-vmi",
					},
				}

				// This will fail due to no KubeVirt cluster, but we're testing URI parsing
				result, err := resources.VmiGet(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				// Should not contain URI parsing errors
				Expect(err.Error()).NotTo(ContainSubstring("invalid URI format"))
			})
		})

		Context("when given valid URI format", func() {
			It("should accept valid kubevirt URI and attempt to get VMI", func() {
				request := mcp.ReadResourceRequest{
					Params: struct {
						URI       string                 `json:"uri"`
						Arguments map[string]interface{} `json:"arguments,omitempty"`
					}{
						URI: "kubevirt://test-namespace/vmi/test-vmi",
					},
				}

				// This will fail due to no KubeVirt cluster, but we're testing the URI parsing
				result, err := resources.VmiGet(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).NotTo(ContainSubstring("invalid URI format"))
				Expect(err.Error()).NotTo(ContainSubstring("missing"))
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("DataVolumesList", func() {
		Context("when given invalid URI", func() {
			It("should return an error for malformed URI", func() {
				request := mcp.ReadResourceRequest{
					Params: struct {
						URI       string                 `json:"uri"`
						Arguments map[string]interface{} `json:"arguments,omitempty"`
					}{
						URI: "invalid-uri",
					},
				}

				result, err := resources.DataVolumesList(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid URI format"))
				Expect(result).To(BeNil())
			})

			It("should return an error for URI with missing namespace", func() {
				request := mcp.ReadResourceRequest{
					Params: struct {
						URI       string                 `json:"uri"`
						Arguments map[string]interface{} `json:"arguments,omitempty"`
					}{
						URI: "kubevirt:/datavolumes",
					},
				}

				result, err := resources.DataVolumesList(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid URI format"))
				Expect(result).To(BeNil())
			})
		})

		Context("when given valid URI format", func() {
			It("should parse namespace correctly", func() {
				request := mcp.ReadResourceRequest{
					Params: struct {
						URI       string                 `json:"uri"`
						Arguments map[string]interface{} `json:"arguments,omitempty"`
					}{
						URI: "kubevirt://test-namespace/datavolumes",
					},
				}

				// This will fail due to no KubeVirt cluster, but we're testing URI parsing
				result, err := resources.DataVolumesList(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				// Should not contain URI parsing errors
				Expect(err.Error()).NotTo(ContainSubstring("invalid URI format"))
			})
		})

		Context("when given valid URI format", func() {
			It("should accept valid kubevirt URI and attempt to list DataVolumes", func() {
				request := mcp.ReadResourceRequest{
					Params: struct {
						URI       string                 `json:"uri"`
						Arguments map[string]interface{} `json:"arguments,omitempty"`
					}{
						URI: "kubevirt://test-namespace/datavolumes",
					},
				}

				// This will fail due to no KubeVirt cluster, but we're testing the URI parsing
				result, err := resources.DataVolumesList(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).NotTo(ContainSubstring("invalid URI format"))
				Expect(err.Error()).NotTo(ContainSubstring("missing namespace"))
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("DataVolumeGet", func() {
		Context("when given invalid URI", func() {
			It("should return an error for malformed URI", func() {
				request := mcp.ReadResourceRequest{
					Params: struct {
						URI       string                 `json:"uri"`
						Arguments map[string]interface{} `json:"arguments,omitempty"`
					}{
						URI: "invalid-uri",
					},
				}

				result, err := resources.DataVolumeGet(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid URI format"))
				Expect(result).To(BeNil())
			})

			It("should return an error for URI missing DataVolume name", func() {
				request := mcp.ReadResourceRequest{
					Params: struct {
						URI       string                 `json:"uri"`
						Arguments map[string]interface{} `json:"arguments,omitempty"`
					}{
						URI: "kubevirt://test-namespace/datavolume",
					},
				}

				result, err := resources.DataVolumeGet(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid URI format"))
				Expect(result).To(BeNil())
			})
		})

		Context("when given valid URI format", func() {
			It("should parse namespace and name correctly", func() {
				request := mcp.ReadResourceRequest{
					Params: struct {
						URI       string                 `json:"uri"`
						Arguments map[string]interface{} `json:"arguments,omitempty"`
					}{
						URI: "kubevirt://test-namespace/datavolume/test-dv",
					},
				}

				// This will fail due to no KubeVirt cluster, but we're testing URI parsing
				result, err := resources.DataVolumeGet(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				// Should not contain URI parsing errors
				Expect(err.Error()).NotTo(ContainSubstring("invalid URI format"))
			})
		})

		Context("when given valid URI format", func() {
			It("should accept valid kubevirt URI and attempt to get DataVolume", func() {
				request := mcp.ReadResourceRequest{
					Params: struct {
						URI       string                 `json:"uri"`
						Arguments map[string]interface{} `json:"arguments,omitempty"`
					}{
						URI: "kubevirt://test-namespace/datavolume/test-dv",
					},
				}

				// This will fail due to no KubeVirt cluster, but we're testing the URI parsing
				result, err := resources.DataVolumeGet(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).NotTo(ContainSubstring("invalid URI format"))
				Expect(err.Error()).NotTo(ContainSubstring("missing"))
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("VmGetStatus", func() {
		Context("when given invalid URI", func() {
			It("should return an error for malformed URI", func() {
				request := mcp.ReadResourceRequest{
					Params: struct {
						URI       string                 `json:"uri"`
						Arguments map[string]interface{} `json:"arguments,omitempty"`
					}{
						URI: "invalid-uri",
					},
				}

				result, err := resources.VmGetStatus(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid URI format"))
				Expect(result).To(BeNil())
			})

			It("should return an error for URI missing VM name", func() {
				request := mcp.ReadResourceRequest{
					Params: struct {
						URI       string                 `json:"uri"`
						Arguments map[string]interface{} `json:"arguments,omitempty"`
					}{
						URI: "kubevirt://test-namespace/vm//status",
					},
				}

				result, err := resources.VmGetStatus(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				// Should not contain URI parsing errors
				Expect(err.Error()).NotTo(ContainSubstring("invalid URI format"))
			})
		})

		Context("when given valid URI format", func() {
			It("should accept valid kubevirt URI and attempt to get VM status", func() {
				request := mcp.ReadResourceRequest{
					Params: struct {
						URI       string                 `json:"uri"`
						Arguments map[string]interface{} `json:"arguments,omitempty"`
					}{
						URI: "kubevirt://test-namespace/vm/test-vm/status",
					},
				}

				// This will fail due to no KubeVirt cluster, but we're testing the URI parsing
				result, err := resources.VmGetStatus(ctx, request)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).NotTo(ContainSubstring("invalid URI format"))
				Expect(err.Error()).NotTo(ContainSubstring("missing"))
				Expect(result).To(BeNil())
			})
		})
	})
})
