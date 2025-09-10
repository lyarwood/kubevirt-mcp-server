package containerdisks_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestContainerDisks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ContainerDisks Suite")
}
