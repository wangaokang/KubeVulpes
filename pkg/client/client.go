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

package client

import (
	"encoding/base64"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func ParseKubeConfigBytes(cfg string) ([]byte, error) {
	kubeConfigBytes, err := base64.StdEncoding.DecodeString(cfg)
	if err != nil {
		return nil, err
	}

	return kubeConfigBytes, err
}

func NewClientSetFromBytes(data []byte) (*kubernetes.Clientset, error) {
	config, err := clientcmd.RESTConfigFromKubeConfig(data)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func NewClientSetFromString(cfg string) (*kubernetes.Clientset, error) {
	kubeConfigBytes, err := ParseKubeConfigBytes(cfg)
	if err != nil {
		return nil, err
	}

	return NewClientSetFromBytes(kubeConfigBytes)
}

func NewClusterSet(cfg string) (*ClusterSet, error) {
	kubeConfigBytes, err := ParseKubeConfigBytes(cfg)
	if err != nil {
		return nil, err
	}

	cs := &ClusterSet{}
	if err = cs.Complete(kubeConfigBytes); err != nil {
		return nil, err
	}

	return cs, nil
}
