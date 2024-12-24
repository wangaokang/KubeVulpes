package cluster

import (
	"context"

	"github.com/casbin/casbin/v2"

	"kubevulpes/cmd/app/config"
	"kubevulpes/pkg/db"
	"kubevulpes/pkg/types"
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

func (c *cluster) Create(ctx context.Context, req *types.CreateClusterRequest) error {
	return nil
}

func (c *cluster) Update(ctx context.Context, clusterId int64, req *types.UpdateClusterRequest) error {
	return nil
}

func (c *cluster) Delete(ctx context.Context, clusterId int64) error {
	return nil
}

func (c *cluster) Get(ctx context.Context, clusterId int64) (*types.Cluster, error) {
	return nil, nil
}

func (c *cluster) List(ctx context.Context, listOptions *types.ListOptions) (*types.PageResponse, error) {
	return nil, nil
}

func New(cc config.Config, f db.ShareDaoFactory, e *casbin.SyncedEnforcer) Interface {
	return &cluster{
		cc:       cc,
		factory:  f,
		enforcer: e,
	}
}
