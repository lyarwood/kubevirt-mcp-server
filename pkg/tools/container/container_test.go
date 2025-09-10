package container_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/lyarwood/kubevirt-mcp-server/pkg/tools/container"
)

var _ = Describe("Container", func() {
	Describe("ResolveContainerDisk", func() {
		Context("when given OS names", func() {
			It("should resolve fedora to containerdisks fedora image", func() {
				result := container.ResolveContainerDisk("fedora")
				Expect(result).To(Equal("quay.io/containerdisks/fedora:latest"))
			})

			It("should resolve ubuntu to containerdisks ubuntu image", func() {
				result := container.ResolveContainerDisk("ubuntu")
				Expect(result).To(Equal("quay.io/containerdisks/ubuntu:latest"))
			})

			It("should resolve centos to containerdisks centos image", func() {
				result := container.ResolveContainerDisk("centos")
				Expect(result).To(Equal("quay.io/containerdisks/centos:latest"))
			})

			It("should resolve cirros to kubevirt demo image", func() {
				result := container.ResolveContainerDisk("cirros")
				Expect(result).To(Equal("quay.io/kubevirt/cirros-container-disk-demo"))
			})

			It("should handle case insensitive input", func() {
				result := container.ResolveContainerDisk("FEDORA")
				Expect(result).To(Equal("quay.io/containerdisks/fedora:latest"))
			})

			It("should handle input with extra whitespace", func() {
				result := container.ResolveContainerDisk("  ubuntu  ")
				Expect(result).To(Equal("quay.io/containerdisks/ubuntu:latest"))
			})
		})

		Context("when given container image names", func() {
			It("should return full container image URLs as-is", func() {
				input := "quay.io/containerdisks/fedora:38"
				result := container.ResolveContainerDisk(input)
				Expect(result).To(Equal(input))
			})

			It("should return images with tags as-is", func() {
				input := "my-registry/my-image:v1.0"
				result := container.ResolveContainerDisk(input)
				Expect(result).To(Equal(input))
			})

			It("should return images with slashes as-is", func() {
				input := "docker.io/library/ubuntu"
				result := container.ResolveContainerDisk(input)
				Expect(result).To(Equal(input))
			})
		})

		Context("when given unknown OS names", func() {
			It("should construct containerdisks URL for unknown OS", func() {
				result := container.ResolveContainerDisk("myos")
				Expect(result).To(Equal("quay.io/containerdisks/myos:latest"))
			})

			It("should handle unknown OS with case normalization", func() {
				result := container.ResolveContainerDisk("MyCustomOS")
				Expect(result).To(Equal("quay.io/containerdisks/mycustomos:latest"))
			})
		})
	})
})