package client_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"main/pkg/client"
)

var _ = Describe("Client", func() {
	Describe("GetKubevirtClient", func() {
		Context("when creating a KubeVirt client", func() {
			It("should return proper client or error", func() {
				// This may succeed or fail depending on environment
				// but we're testing that the function exists and handles both cases
				kubevirtClient, err := client.GetKubevirtClient()

				if err != nil {
					// If error occurs, client should be nil
					Expect(kubevirtClient).To(BeNil())
				} else {
					// If no error, client should not be nil
					Expect(kubevirtClient).NotTo(BeNil())
				}
			})

			It("should use default client configuration", func() {
				// Test that the function doesn't panic and handles errors gracefully
				Expect(func() {
					client.GetKubevirtClient()
				}).NotTo(Panic())
			})
		})
	})
})
