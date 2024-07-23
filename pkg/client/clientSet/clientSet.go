package clientSet

import (
	"k8s.io/client-go/kubernetes"
	"kubevulpes/pkg/client"
)

// EMR测试连通性新建的clientSet
func NewClientSet(data []byte) (*kubernetes.Clientset, error) {
	kubeConfig, err := client.BuildClientConfig(data)
	if err != nil {
		return nil, err
	}
	print()
	client, err := kubernetes.NewForConfig(kubeConfig)

	return client, nil
}
