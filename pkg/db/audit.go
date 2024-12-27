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
	"gorm.io/gorm"
	"kubevulpes/pkg/db/model"
)

type AuditInterface interface {
	Create(ctx context.Context, object *model.Audit) error
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (*model.Audit, error)
	List(ctx context.Context, opts ...Options) ([]model.Audit, int64, error)
}

type audit struct {
	db *gorm.DB
}

func (a *audit) Create(ctx context.Context, object *model.Audit) error {
	return a.db.WithContext(ctx).Create(object).Error
}

func (a *audit) Delete(ctx context.Context, id int64) error {
	return a.db.WithContext(ctx).Delete(&model.Audit{}, id).Error
}

func (a *audit) Get(ctx context.Context, id int64) (*model.Audit, error) {
	var au model.Audit
	if err := a.db.WithContext(ctx).First(&au, id).Error; err != nil {
		return nil, err
	}
	return &au, nil
}

func (a *audit) List(ctx context.Context, opts ...Options) ([]model.Audit, int64, error) {
	var (
		audits []model.Audit
		total  int64
	)

	if err := a.db.WithContext(ctx).Model(&model.Audit{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := a.db.WithContext(ctx).Find(&audits).Error; err != nil {
		return nil, 0, err
	}
	return audits, total, nil
}

func newAudit(db *gorm.DB) *audit {
	return &audit{db}
}
