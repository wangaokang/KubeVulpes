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

package controller

import (
	"github.com/casbin/casbin/v2"
	"kubevulpes/pkg/controller/audit"
	"kubevulpes/pkg/controller/auth"

	"kubevulpes/cmd/app/config"
	"kubevulpes/pkg/controller/cluster"
	"kubevulpes/pkg/controller/user"
	"kubevulpes/pkg/db"
)

type VuplesInterface interface {
	user.UserGetter
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
