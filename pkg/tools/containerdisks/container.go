package containerdisks

import (
	"fmt"
	"strings"
)

// ResolveContainerDisk resolves OS names to container disk images from quay.io/containerdisks
func ResolveContainerDisk(input string) string {
	// If input already looks like a container image, return as-is
	if strings.Contains(input, "/") || strings.Contains(input, ":") {
		return input
	}

	// Common OS name mappings to containerdisk images
	osMap := map[string]string{
		"fedora":   "quay.io/containerdisks/fedora:latest",
		"ubuntu":   "quay.io/containerdisks/ubuntu:latest",
		"centos":   "quay.io/containerdisks/centos:latest",
		"debian":   "quay.io/containerdisks/debian:latest",
		"rhel":     "quay.io/containerdisks/rhel:latest",
		"opensuse": "quay.io/containerdisks/opensuse:latest",
		"alpine":   "quay.io/containerdisks/alpine:latest",
		"cirros":   "quay.io/kubevirt/cirros-container-disk-demo",
		"windows":  "quay.io/containerdisks/windows:latest",
		"freebsd":  "quay.io/containerdisks/freebsd:latest",
	}

	// Normalize input to lowercase for lookup
	normalized := strings.ToLower(strings.TrimSpace(input))

	// Look up the OS name
	if containerDisk, exists := osMap[normalized]; exists {
		return containerDisk
	}

	// If no match found, assume it's already a valid container disk name
	// and try to construct a containerdisks URL
	return fmt.Sprintf("quay.io/containerdisks/%s:latest", normalized)
}
