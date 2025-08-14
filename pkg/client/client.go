package client

import (
	"github.com/spf13/pflag"
	"kubevirt.io/client-go/kubecli"
)

// GetKubevirtClient returns a KubeVirt client using the default configuration
func GetKubevirtClient() (kubecli.KubevirtClient, error) {
	clientConfig := kubecli.DefaultClientConfig(&pflag.FlagSet{})
	return kubecli.GetKubevirtClientFromClientConfig(clientConfig)
}