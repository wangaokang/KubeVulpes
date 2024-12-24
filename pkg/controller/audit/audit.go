package audit

import (
	"context"
	"github.com/casbin/casbin/v2"
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
	//object, err := a.factory.Audit().Get(ctx, aid)
	//if err != nil {
	//	klog.Errorf("failed to get audit %d: %v", aid, err)
	//	return nil, errors.ErrServerInternal
	//}
	//if object == nil {
	//	return nil, errors.ErrAuditNotFound
	//}
	//
	//return a.model2Type(object), nil
	return nil, nil
}

func (a *audit) List(ctx context.Context, listOption types.ListOptions) (interface{}, error) {
	//var ts []types.Audit
	//
	//
	//// 获取偏移列表
	//objects, total,err := a.factory.Audit().List(ctx, listOption.BuildPageNation()...)
	//if err != nil {
	//	klog.Errorf("failed to get audit events: %v", err)
	//	return nil, err
	//}
	//
	//for _, object := range objects {
	//	ts = append(ts, *a.model2Type(&object))
	//}
	//
	//return types.PageResponse{
	//	PageRequest: listOption.PageRequest,
	//	Total:       int(total),
	//	Items:       ts,
	//}, nil

	return nil, nil
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
		Ip:         o.Ip,
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
