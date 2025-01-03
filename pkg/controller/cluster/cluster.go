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

package cluster

import (
	"context"
	"fmt"
	"net/http"

	"github.com/casbin/casbin/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"

	"kubevulpes/api/errors"
	"kubevulpes/cmd/app/config"
	"kubevulpes/pkg/client"
	ctrlutil "kubevulpes/pkg/controller/util"
	"kubevulpes/pkg/db"
	"kubevulpes/pkg/db/model"
	"kubevulpes/pkg/types"
	"kubevulpes/pkg/util/uuid"
)

type ClusterGetter interface {
	Cluster() Interface
}

type Interface interface {
	Create(ctx context.Context, req *types.CreateClusterRequest) error
	Update(ctx context.Context, clusterId int64, req *types.UpdateClusterRequest) error
	Delete(ctx context.Context, clusterId int64) error
	Get(ctx context.Context, clusterId int64) (*types.Cluster, error)
	List(ctx context.Context, listOptions *types.ListOptions) (*types.PageResponse, error)
}

type cluster struct {
	cc       config.Config
	factory  db.ShareDaoFactory
	enforcer *casbin.SyncedEnforcer
}

var clusterIndexer client.Cache

func init() {
	clusterIndexer = *client.NewClusterCache()
}

func (c *cluster) Create(ctx context.Context, req *types.CreateClusterRequest) error {
	//user, err := httputils.GetUserFromRequest(ctx)
	//if err != nil {
	//	return errors.NewError(err, http.StatusInternalServerError)
	//}

	if err := c.preCreate(ctx, req); err != nil {
		return errors.NewError(err, http.StatusBadRequest)
	}
	// TODO: 集群名称必须是由英文，数字组成
	if len(req.Name) == 0 {
		req.Name = uuid.NewRandName(8)
	}

	var cs *client.ClusterSet
	var txFunc = func(cluster *model.Cluster) (err error) {
		if cs, err = client.NewClusterSet(req.KubeConfig); err != nil {
			return
		}

		// insert a user RBAC policy 后续考虑添加集群权限
		//policy := model.NewPolicyFromModels(user, model.ObjectCluster, cluster.Model, model.OpAll)
		//_, err = c.enforcer.AddPolicy(policy.Raw())
		return
	}

	if _, err := c.factory.Cluster().Create(ctx, &model.Cluster{
		Name:       req.Name,
		KubeConfig: req.KubeConfig,
	}, txFunc); err != nil {
		return errors.NewError(err, http.StatusInternalServerError)
	}

	clusterIndexer.Set(req.Name, *cs)
	return nil
}

func (c *cluster) Update(ctx context.Context, cid int64, req *types.UpdateClusterRequest) error {
	object, err := c.factory.Cluster().Get(ctx, cid)
	if err != nil {
		klog.Errorf("failed to get cluster(%d): %v", cid, err)
		return errors.ErrServerInternal
	}
	if object == nil {
		return errors.ErrClusterNotFound
	}
	updates := make(map[string]interface{})
	if req.AliasName != nil {
		updates["alias_name"] = *req.AliasName
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if len(updates) == 0 {
		return errors.ErrInvalidRequest
	}
	if err := c.factory.Cluster().Update(ctx, cid, *req.ResourceVersion, updates); err != nil {
		klog.Errorf("failed to update cluster(%d): %v", cid, err)
		return errors.ErrServerInternal
	}

	return nil
}

func (c *cluster) Delete(ctx context.Context, cid int64) error {
	//user, err := httputils.GetUserFromRequest(ctx)
	//if err != nil {
	//	return errors.NewError(err, http.StatusInternalServerError)
	//}

	cluster, err := c.preDelete(ctx, cid)
	if err != nil {
		return err
	}

	// 后续考虑集群添加权限问题，暂时先注释掉
	var txFunc = func(cluster *model.Cluster) (err error) {
		//_, err = c.enforcer.RemoveNamedPolicy("p", user.Name, model.ObjectCluster.String(), cluster.GetSID())
		return
	}
	if err = c.factory.Cluster().Delete(ctx, cluster, txFunc); err != nil {
		klog.Errorf("failed to delete cluster(%d): %v", cid, err)
		return errors.ErrServerInternal
	}

	// 从缓存中移除 clusterSet
	clusterIndexer.Delete(cluster.Name)
	return nil
}

func (c *cluster) Get(ctx context.Context, cid int64) (*types.Cluster, error) {
	object, err := c.factory.Cluster().Get(ctx, cid)
	if err != nil {
		return nil, errors.ErrServerInternal
	}
	if object == nil {
		return nil, errors.ErrClusterNotFound
	}

	return c.model2Type(object), nil
}

func (c *cluster) List(ctx context.Context, listOptions *types.ListOptions) (*types.PageResponse, error) {
	opts := append(ctrlutil.MakeDbOptions(ctx), listOptions.BuildPageNation()...)

	objects, total, err := c.factory.Cluster().List(ctx, opts...)
	if err != nil {
		return nil, err
	}

	cs := make([]types.Cluster, len(objects))
	for i, object := range objects {
		cs[i] = *c.model2Type(&object)
	}

	return &types.PageResponse{
		Total:       int(total),
		Items:       cs,
		PageRequest: listOptions.PageRequest,
	}, nil
}

func (c *cluster) preCreate(ctx context.Context, req *types.CreateClusterRequest) error {
	// 实际创建前，先创建集群的连通性
	if err := c.Ping(ctx, req.KubeConfig); err != nil {
		return fmt.Errorf("尝试连接 kubernetes API 失败: %v", err)
	}
	return nil
}

// 删除前置检查
// 开启集群删除保护，则不允许删除
func (c *cluster) preDelete(ctx context.Context, cid int64) (cluster *model.Cluster, err error) {
	if cluster, err = c.factory.Cluster().Get(ctx, cid); err != nil {
		klog.Errorf("failed to get cluster(%d): %v", cid, err)
		return
	}
	if cluster == nil {
		return nil, errors.ErrClusterNotFound
	}
	// 开启集群删除保护，则不允许删除
	if cluster.Protected {
		return nil, errors.NewError(fmt.Errorf("已开启集群删除保护功能，不允许删除 %s", cluster.Name),
			http.StatusForbidden)
	}

	return
}

func (c *cluster) model2Type(o *model.Cluster) *types.Cluster {
	return &types.Cluster{
		VulpesMeta: types.VulpesMeta{
			Id:              o.Id,
			ResourceVersion: o.ResourceVersion,
		},
		TimeMeta: types.TimeMeta{
			GmtCreate:   o.GmtCreate,
			GmtModified: o.GmtModified,
		},
		Name:              o.Name,
		AliasName:         o.AliasName,
		KubernetesVersion: o.KubernetesVersion,
		Status:            o.ClusterStatus, // 默认是运行中状态
		Protected:         o.Protected,
		Description:       o.Description,
	}
}

// Ping 检查和 k8s 集群的连通性
// 如果能获取到 k8s 接口的正常返回，则返回 nil，否则返回具体 error
// kubeConfig 为 k8s 证书的 base64 字符串
func (c *cluster) Ping(ctx context.Context, kubeConfig string) error {
	clientSet, err := client.NewClientSetFromString(kubeConfig)
	if err != nil {
		return err
	}

	// 调用 ns 资源，确保连通
	var timeout int64 = 1
	if _, err = clientSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{
		TimeoutSeconds: &timeout,
	}); err != nil {
		klog.Errorf("failed to ping kubernetes: %v", err)
		// 处理原始报错信息，仅返回连接不通的信息
		return fmt.Errorf("kubernetes 集群连接测试失败")
	}

	return nil
}

func NewCluster(cc config.Config, f db.ShareDaoFactory, e *casbin.SyncedEnforcer) Interface {
	return &cluster{
		cc:       cc,
		factory:  f,
		enforcer: e,
	}
}
