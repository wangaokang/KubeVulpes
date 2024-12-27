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

package audit

import (
	"context"
	"github.com/casbin/casbin/v2"
	"k8s.io/klog/v2"
	"kubevulpes/api/errors"
	"kubevulpes/cmd/app/config"
	"kubevulpes/pkg/db"
	"kubevulpes/pkg/db/model"
	"kubevulpes/pkg/types"
)

type AuditGetter interface {
	Audit() Interface
}

type Interface interface {
	Get(ctx context.Context, aid int64) (*types.Audit, error)
	List(ctx context.Context, listOption types.ListOptions) (interface{}, error)
}

type audit struct {
	cc       config.Config
	factory  db.ShareDaoFactory
	enforcer *casbin.SyncedEnforcer
}

func (a *audit) Get(ctx context.Context, aid int64) (*types.Audit, error) {
	object, err := a.factory.Audit().Get(ctx, aid)
	if err != nil {
		klog.Errorf("failed to get audit %d: %v", aid, err)
		return nil, errors.ErrServerInternal
	}
	if object == nil {
		return nil, errors.ErrAuditNotFound
	}

	return a.model2Type(object), nil
}

func (a *audit) List(ctx context.Context, listOption types.ListOptions) (interface{}, error) {
	var ts []types.Audit

	// 获取偏移列表
	objects, total, err := a.factory.Audit().List(ctx, listOption.BuildPageNation()...)
	if err != nil {
		klog.Errorf("failed to get audit events: %v", err)
		return nil, err
	}

	for _, object := range objects {
		ts = append(ts, *a.model2Type(&object))
	}

	return types.PageResponse{
		PageRequest: listOption.PageRequest,
		Total:       int(total),
		Items:       ts,
	}, nil
}

func (a *audit) model2Type(o *model.Audit) *types.Audit {
	return &types.Audit{
		PixiuMeta: types.PixiuMeta{
			Id:              o.Id,
			ResourceVersion: o.ResourceVersion,
		},
		TimeMeta: types.TimeMeta{
			GmtCreate:   o.GmtCreate,
			GmtModified: o.GmtModified,
		},
		Ip:         o.IP,
		Action:     o.Action,
		Status:     o.Status,
		Operator:   o.Operator,
		Path:       o.Path,
		ObjectType: o.ObjectType,
	}
}

func NewAudit(cfg config.Config, f db.ShareDaoFactory) *audit {
	return &audit{
		cc:      cfg,
		factory: f,
	}
}
