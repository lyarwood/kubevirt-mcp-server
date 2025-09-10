package preference_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPreference(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Preference Suite")
}
