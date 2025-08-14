package functional_test

import (
	"context"
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/tools/clientcmd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/api/core/v1"
	"kubevirt.io/client-go/kubecli"
)

var (
	virtClient   kubecli.KubevirtClient
	testTimeout  = 300 * time.Second
	pollInterval = 5 * time.Second
)

func TestFunctional(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Functional Test Suite")
}

var _ = BeforeSuite(func() {
	// Get kubeconfig path from environment or use kubevirtci default
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		kubeconfigPath = "../../_kubevirtci/_ci-configs/k8s-1.33/.kubeconfig"
	}

	// Check if kubeconfig exists
	if _, err := os.Stat(kubeconfigPath); os.IsNotExist(err) {
		Skip("Kubeconfig not found. Please run 'make cluster-up' first")
	}

	// Create Kubernetes client config
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	Expect(err).NotTo(HaveOccurred())

	// Create KubeVirt client
	virtClient, err = kubecli.GetKubevirtClientFromRESTConfig(config)
	Expect(err).NotTo(HaveOccurred())

	// Verify KubeVirt is deployed and ready
	By("Verifying KubeVirt is ready")
	Eventually(func() bool {
		kv, err := virtClient.KubeVirt("kubevirt").Get(context.Background(), "kubevirt", metav1.GetOptions{})
		if err != nil {
			return false
		}
		return kv.Status.Phase == kubevirtv1.KubeVirtPhaseDeployed
	}, testTimeout, pollInterval).Should(BeTrue(), "KubeVirt should be deployed and ready")
})

var _ = AfterSuite(func() {
	// Cleanup can be added here if needed
})