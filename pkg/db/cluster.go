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

type ClusterInterface interface {
	Create(ctx context.Context, object *model.Cluster) error
	Delete(ctx context.Context, clusterId int64) error
	Get(ctx context.Context, clusterId int64) (*model.Cluster, error)
	List(ctx context.Context, opts ...Options) ([]model.Cluster, int64, error)
	Update(ctx context.Context, clusterId int64, resourceVersion int64, updates map[string]interface{}) error
}

type cluster struct {
	db *gorm.DB
}

func (c *cluster) Create(ctx context.Context, object *model.Cluster) error {
	return c.db.WithContext(ctx).Create(object).Error
}

func (c *cluster) Delete(ctx context.Context, clusterId int64) error {
	return c.db.WithContext(ctx).Delete(&model.Cluster{}, clusterId).Error
}

func (c *cluster) Get(ctx context.Context, clusterId int64) (*model.Cluster, error) {
	var cluster model.Cluster
	if err := c.db.WithContext(ctx).First(&cluster, clusterId).Error; err != nil {
		return nil, err
	}
	return &cluster, nil
}

func (c *cluster) List(ctx context.Context, opts ...Options) ([]model.Cluster, int64, error) {
	var (
		clusters []model.Cluster
		total    int64
	)

	if err := c.db.WithContext(ctx).Find(&clusters).Error; err != nil {
		return nil, 0, err
	}
	if err := c.db.WithContext(ctx).Model(&model.Cluster{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return clusters, total, nil
}

func (c *cluster) Update(ctx context.Context, clusterId int64, resourceVersion int64, updates map[string]interface{}) error {
	return c.db.WithContext(ctx).Model(&model.Cluster{}).Where("id = ? and resource_version = ?", clusterId, resourceVersion).Updates(updates).Error
}

func newCluster(db *gorm.DB) ClusterInterface {
	return &cluster{db: db}
}
