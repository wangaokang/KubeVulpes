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
