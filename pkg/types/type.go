package types

import (
	"io"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
	appv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/remotecommand"

	"kubevulpes/pkg/db/model"
)

type PixiuObjectMeta struct {
	Cluster   string `uri:"cluster" binding:"required"`
	Namespace string `uri:"namespace" binding:"required"`
	Name      string `uri:"name"`
}

type PixiuMeta struct {
	// pixiu 对象 ID
	Id int64 `json:"id"`
	// Pixiu 对象版本号
	ResourceVersion int64 `json:"resource_version"`
}

type TimeMeta struct {
	// pixiu 对象创建时间
	GmtCreate time.Time `json:"gmt_create"`
	// pixiu 对象修改时间
	GmtModified time.Time `json:"gmt_modified"`
}

type KubeNode struct {
	Ready    []string `json:"ready"`
	NotReady []string `json:"not_ready"`
}

type Cluster struct {
	PixiuMeta `json:",inline"`

	Name      string `json:"name"`
	AliasName string `json:"alias_name"`
	//Status    model.ClusterStatus `json:"status"` // 0: 运行中 1: 部署中 2: 等待部署 3: 部署失败 4: 集群失联，API不可用

	// 0: 标准集群 1: 自建集群
	//ClusterType model.ClusterType `json:"cluster_type"`
	PlanId int64 `json:"plan_id"` // 自建集群关联的 PlanId，如果是自建的集群，planId 不为 0

	// kubernetes 集群的版本和状态
	KubernetesVersion string   `json:"kubernetes_version"`
	Nodes             KubeNode `json:"nodes"`

	// 集群删除保护，开启集群删除保护时不允许删除集群
	// 0: 关闭集群删除保护 1: 开启集群删除保护
	Protected bool `json:"protected"`

	// k8s kubeConfig base64 字段
	KubeConfig string `json:"kube_config,omitempty"`

	// 集群用途描述，可以为空
	Description string `json:"description"`

	KubernetesMeta `json:",inline"`
	TimeMeta       `json:",inline"`
}

// KubernetesMeta 记录 kubernetes 集群的数据
type KubernetesMeta struct {
	// 集群的版本
	KubernetesVersion string `json:"kubernetes_version,omitempty"`
	// 节点数量
	Nodes int `json:"nodes"`
	// The memory and cpu usage
	Resources Resources `json:"resources"`
}

// Resources kubernetes 的资源信息
// The memory and cpu usage
type Resources struct {
	Cpu    string `json:"cpu"`
	Memory string `json:"memory"`
}

type User struct {
	PixiuMeta `json:",inline"`

	Name        string           `json:"name"`                                 // 用户名称
	Password    string           `json:"password" binding:"required,password"` // 用户密码
	Status      model.UserStatus `json:"status"`                               // 用户状态标识
	Role        model.UserRole   `json:"role"`                                 // 用户角色，目前只实现管理员，0: 普通用户 1: 管理员 2: 超级管理员
	Email       string           `json:"email"`                                // 用户注册邮件
	Description string           `json:"description"`                          // 用户描述信息

	TimeMeta `json:",inline"`
}

type Role struct {
	PixiuMeta `json:",inline"`

	Name string `json:"name"`
}

type Tenant struct {
	PixiuMeta `json:",inline"`
	TimeMeta  `json:",inline"`

	Name        string `json:"name"`        // 用户名称
	Description string `json:"description"` // 用户描述信息
}

type Audit struct {
	PixiuMeta `json:",inline"`
	TimeMeta  `json:",inline"`

	Ip         string                     `json:"ip"`
	Action     string                     `json:"action"`        // 操作动作
	Status     model.AuditOperationStatus `json:"status"`        // 操作状态
	Operator   string                     `json:"operator"`      // 操作人
	Path       string                     `json:"path"`          // 操作路径
	ObjectType model.ObjectType           `json:"resource_type"` // 资源类型
}

type AuthType string

type KeySpec struct {
	Data string `json:"data,omitempty"`
	File string `json:"-"`
}

type PasswordSpec struct {
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}

// TimeSpec 通用时间规格
type TimeSpec struct {
	GmtCreate   interface{} `json:"gmt_create,omitempty"`
	GmtModified interface{} `json:"gmt_modified,omitempty"`
}

type KubeObject struct {
	lock sync.RWMutex

	ReplicaSets []appv1.ReplicaSet
	Pods        []v1.Pod
}

// WebShellOptions ws API 参数定义
type WebShellOptions struct {
	Cluster   string `form:"cluster"`
	Namespace string `form:"namespace"`
	Pod       string `form:"pod"`
	Container string `form:"container"`
	Command   string `form:"command"`
}

// TerminalMessage 定义了终端和容器 shell 交互内容的格式 Operation 是操作类型
// Data 是具体数据内容 Rows和Cols 可以理解为终端的行数和列数，也就是宽、高
type TerminalMessage struct {
	Operation string `json:"operation"`
	Data      string `json:"data"`
	Rows      uint16 `json:"rows"`
	Cols      uint16 `json:"cols"`
}

// TerminalSession 定义 TerminalSession 结构体，实现 PtyHandler 接口
// wsConn 是 websocket 连接
// sizeChan 用来定义终端输入和输出的宽和高
// doneChan 用于标记退出终端
type TerminalSession struct {
	wsConn   *websocket.Conn
	sizeChan chan remotecommand.TerminalSize
	doneChan chan struct{}
}

type Turn struct {
	StdinPipe io.WriteCloser
	Session   *ssh.Session
	WsConn    *websocket.Conn
}

// ListOptions is the query options to a standard REST list call.
type ListOptions struct {
	Count bool `form:"count"`
	//Limit int64 `form:"limit"`

	PageRequest `json:",inline"` // 分页请求属性
	QueryOption `json:",inline"` // 搜索内容
}

type EventOptions struct {
	Uid        string `form:"uid"`
	Namespace  string `form:"namespace"`
	Name       string `form:"name"`
	Kind       string `form:"kind"`
	Namespaced bool   `form:"namespaced"`
	Limit      int64  `form:"limit"`
}

type PodLogOptions struct {
	Container string `form:"container"`
	TailLines int64  `form:"tailLines"`
}

type KubernetesSpec struct {
	EnablePublicIp    bool   `json:"enable_public_ip"`
	ApiServer         string `json:"api_server"`
	ApiPort           string `json:"api_port"`
	KubernetesVersion string `json:"kubernetes_version"`
	EnableHA          bool   `json:"enable_ha"`
	Register          bool   `json:"register"`
}

type NetworkSpec struct {
	NetworkInterface string `json:"network_interface"` // 网口，默认 eth0
	Cni              string `json:"cni"`
	PodNetwork       string `json:"pod_network"`
	ServiceNetwork   string `json:"service_network"`
	KubeProxy        string `json:"kube_proxy"`
}

type RuntimeSpec struct {
	Runtime string `json:"runtime"`
}

type ComponentSpec struct {
	Helm       *Helm       `json:"helm,omitempty"` // 忽略，则使用默认值
	Prometheus *Prometheus `json:"prometheus,omitempty"`
	Grafana    *Grafana    `json:"grafana,omitempty"`
	Haproxy    *Haproxy    `json:"haproxy,omitempty"`
}

type Helm struct {
	Enable      bool   `json:"enable"`
	HelmRelease string `json:"helm_release"`
}

type Prometheus struct {
	EnablePrometheus string `json:"enable_prometheus"`
	Enable           bool   `json:"enable"`
}

type Grafana struct {
	Enable               bool   `json:"enable"`
	GrafanaAdminUser     string `json:"grafana_admin_user"`
	GrafanaAdminPassword string `json:"grafana_admin_password"`
}

// Haproxy Options
// This configuration is usually enabled when self-created VMs require high availability.
type Haproxy struct {
	Enable                    bool   `json:"enable"`                       // Enable haproxy and keepalived,
	KeepalivedVirtualRouterId string `json:"keepalived_virtual_router_id"` // Arbitrary unique number from 0..255
}

type RBACPolicy struct {
	UserName   string           `json:"username,omitempty"`
	GroupName  string           `json:"groupname,omitempty"`
	ObjectType model.ObjectType `json:"resource_type,omitempty"`
	StringID   string           `json:"sid,omitempty"`
	Operation  model.Operation  `json:"operation,omitempty"`
}
