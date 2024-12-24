package role

import (
	"context"
	"github.com/casbin/casbin/v2"
	"kubevulpes/pkg/types"

	"kubevulpes/cmd/app/config"
	"kubevulpes/pkg/db"
)

type RoleGetter interface {
	Role() Interface
}

type Interface interface {
	Create(ctx context.Context, req *types.CreateRoleRequest) error
	Delete(ctx context.Context, roleId int64) error
	Update(ctx context.Context, roleId int64, req *types.UpdateRoleRequest) error
	Get(ctx context.Context, roleId int64) (*types.Role, error)
	List(ctx context.Context, listOptions *types.ListOptions) (*types.PageResponse, error)
}

type role struct {
	cc       config.Config
	factory  db.ShareDaoFactory
	enforcer *casbin.SyncedEnforcer
}

func (r *role) Create(ctx context.Context, req *types.CreateRoleRequest) error {
	return nil
}

func (r *role) Delete(ctx context.Context, roleId int64) error {
	return nil
}

func (r *role) Get(ctx context.Context, roleId int64) (*types.Role, error) {
	return nil, nil
}

func (r *role) List(ctx context.Context, listOptions *types.ListOptions) (*types.PageResponse, error) {
	return nil, nil
}

func (r *role) Update(ctx context.Context, roleId int64, req *types.UpdateRoleRequest) error {
	return nil
}

func NewRole(cfg config.Config, f db.ShareDaoFactory, e *casbin.SyncedEnforcer) *role {
	return &role{
		cc:       cfg,
		factory:  f,
		enforcer: e,
	}
}
