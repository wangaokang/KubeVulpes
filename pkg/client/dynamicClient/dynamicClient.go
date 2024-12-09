package dynamicClient

import (
	"fmt"

	"k8s.io/client-go/dynamic"

	"kubevulpes/pkg/client"
)

// dynamic.DynamicClient 这个类型在k8s.io/client-go v0.22.1不存在，需要升级到k8s.io/client-go v0.26.3
func NewKubeDynamicClient(Config []byte) (*dynamic.DynamicClient, error) {
	kubeConfig, err := client.BuildClientConfig(Config)
	if err != nil {
		return nil, fmt.Errorf("failed to init kube config: [%v]", err)
	}
	dynamicClient, err := dynamic.NewForConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to NewForConfig client: [%v]", err)
	}

	return dynamicClient, nil
}
