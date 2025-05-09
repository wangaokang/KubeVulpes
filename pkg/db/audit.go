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

package db

import (
	"context"
	"time"

	"gorm.io/gorm"

	"kubevulpes/pkg/db/model"
	"kubevulpes/pkg/util/errors"
)

type AuditInterface interface {
	Create(ctx context.Context, object *model.Audit) error
	BatchDelete(ctx context.Context, opts ...Options) (int64, error)
	Get(ctx context.Context, id int64) (*model.Audit, error)
	List(ctx context.Context, opts ...Options) ([]model.Audit, int64, error)
}

type audit struct {
	db *gorm.DB
}

func (a *audit) Create(ctx context.Context, object *model.Audit) error {
	now := time.Now()
	object.GmtCreate = now
	object.GmtModified = now

	if err := a.db.WithContext(ctx).Create(object).Error; err != nil {
		return err
	}
	return nil
}

func (a *audit) BatchDelete(ctx context.Context, opts ...Options) (int64, error) {
	tx := a.db.WithContext(ctx)
	for _, opt := range opts {
		tx = opt(tx)
	}

	err := tx.Delete(&model.Audit{}).Error
	return tx.RowsAffected, err
}

func (a *audit) Get(ctx context.Context, aid int64) (*model.Audit, error) {
	var audit *model.Audit
	if err := a.db.WithContext(ctx).Where("id = ?", aid).First(audit).Error; err != nil {
		if errors.IsRecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return audit, nil
}

func (a *audit) List(ctx context.Context, opts ...Options) ([]model.Audit, int64, error) {
	var (
		audits []model.Audit
		total  int64
	)

	tx := a.db.WithContext(ctx)
	for _, opt := range opts {
		tx = opt(tx)
	}

	if err := tx.Model(&model.Audit{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := tx.Find(&audits).Error; err != nil {
		return nil, 0, err
	}

	return audits, total, nil
}

func newAudit(db *gorm.DB) *audit {
	return &audit{db}
}
