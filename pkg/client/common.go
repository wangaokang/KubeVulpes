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
	"path/filepath"

	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

//// 测试集群连通性
//func Ping(ctx context.Context, kubeConfigData []byte, ns string) error {
//	clientSet, err := clientSet.NewClientSet(kubeConfigData)
//	if err != nil {
//		return fmt.Errorf("failed to create clientSet for Ping: [%v]", err)
//	}
//	if ns == "" {
//		ns = "default"
//	}
//	_, err = clientSet.AppsV1().Deployments(ns).List(ctx, metav1.ListOptions{})
//	if err != nil {
//		return fmt.Errorf("failed to create clientSet for Ping: [%v]", err)
//	}
//	return nil
//}

// 封装一个使用kubeconfig的逻辑处理
func BuildClientConfig(configFile []byte) (*restclient.Config, error) {
	//todo 后续通过cobra 获取参数
	if len(configFile) == 0 {
		configFile = []byte(filepath.Join(homedir.HomeDir(), ".kube", "config"))
	}

	config, err := clientcmd.RESTConfigFromKubeConfig(configFile)
	if err != nil {
		return nil, err
	}

	return config, nil
}
