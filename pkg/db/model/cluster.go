package model

import "kubevulpes/pkg/db/model/base"

type ClusterStatus uint8

const (
	ClusterStatusRunning ClusterStatus = iota // 运行中
	ClusterStatusDeploy                       // 部署中
	ClusterStatusUnStart                      // 等待部署
	ClusterStatusFailed                       // 部署失败
	ClusterStatusError                        // 集群失联，API不可用
)

func init() {
	register(&Cluster{})
}

type Cluster struct {
	base.Model
	Name string `gorm:"column:name;types:varchar(128);not null" json:"name"`
	// 集群运行状态 0: 运行中 2: 集群失联 4: 所有的 node 不健康
	ClusterStatus `gorm:"column:status;types:tinyint;not null" json:"status"`
	// 集群的版本
	KubernetesVersion string `gorm:"type:varchar(255)" json:"kubernetes_version,omitempty"`
	// 集群节点健康数，json 字符串
	Nodes string `gorm:"type:text" json:"nodes"`
	// 集群删除保护，开启集群删除保护时不允许删除集群
	// 0: 关闭集群删除保护 1: 开启集群删除保护
	Protected bool `json:"protected"`
	// k8s kubeConfig todo 后续考虑对kubeconfig 进行加密
	KubeConfig string `json:"kube_config"`
	// 集群用途描述，可以为空
	Description string `gorm:"type:text" json:"description"`
}
