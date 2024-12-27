/*
Copyright 2024 The Vuples Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

	c, err := kubernetes.NewForConfig(kubeConfig)

	return c, nil
}
