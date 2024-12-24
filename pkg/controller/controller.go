package controller

import (
	"github.com/casbin/casbin/v2"
	"kubevulpes/pkg/controller/audit"
	"kubevulpes/pkg/controller/auth"

	"kubevulpes/cmd/app/config"
	"kubevulpes/pkg/controller/cluster"
	"kubevulpes/pkg/controller/role"
	"kubevulpes/pkg/controller/user"
	"kubevulpes/pkg/db"
)

type VuplesInterface interface {
	user.UserGetter
	role.RoleGetter
	cluster.ClusterGetter
	auth.AuthGetter
	audit.AuditGetter
}

type vuples struct {
	cc       config.Config
	factory  db.ShareDaoFactory
	enforcer *casbin.SyncedEnforcer
}

func (p *vuples) User() user.Interface       { return user.NewUser(p.cc, p.factory, p.enforcer) }
func (p *vuples) Role() role.Interface       { return role.NewRole(p.cc, p.factory, p.enforcer) }
func (p *vuples) Cluster() cluster.Interface { return cluster.New(p.cc, p.factory, p.enforcer) }
func (p *vuples) Auth() auth.Interface       { return auth.NewAuth(p.factory, p.enforcer) }
func (p *vuples) Audit() audit.Interface     { return audit.NewAudit(p.cc, p.factory) }

func New(cfg config.Config, f db.ShareDaoFactory, e *casbin.SyncedEnforcer) VuplesInterface {
	return &vuples{
		cc:       cfg,
		factory:  f,
		enforcer: e,
	}
}
