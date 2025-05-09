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
	"github.com/gin-gonic/gin"

	"kubevulpes/api/httputils"
	"kubevulpes/pkg/types"
)

type IdMeta struct {
	ClusterId int64 `uri:"clusterId" binding:"required"`
}

func (cr *clusterRouter) createCluster(c *gin.Context) {
	r := httputils.NewResponse()
	var (
		req types.CreateClusterRequest
		err error
	)
	if err = c.ShouldBindJSON(&req); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	if err = cr.c.Cluster().Create(c, &req); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	httputils.SetSuccess(c, r)
}

func (cr *clusterRouter) updateCluster(c *gin.Context) {
	r := httputils.NewResponse()
	var (
		idMeta IdMeta
		err    error
	)

	if err = c.ShouldBindUri(&idMeta); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	var req types.UpdateClusterRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	if err = cr.c.Cluster().Update(c, idMeta.ClusterId, &req); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	httputils.SetSuccess(c, r)
}

func (cr *clusterRouter) deleteCluster(c *gin.Context) {
	r := httputils.NewResponse()

	var (
		idMeta IdMeta
		err    error
	)

	if err = c.ShouldBindUri(&idMeta); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	if err = cr.c.Cluster().Delete(c, idMeta.ClusterId); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	httputils.SetSuccess(c, r)
}

func (cr *clusterRouter) getCluster(c *gin.Context) {
	r := httputils.NewResponse()

	var (
		idMeta IdMeta
		err    error
	)

	if err = c.ShouldBindUri(&idMeta); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	if r.Result, err = cr.c.Cluster().Get(c, idMeta.ClusterId); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	httputils.SetSuccess(c, r)
}

func (cr *clusterRouter) listCluster(c *gin.Context) {
	r := httputils.NewResponse()

	var (
		err         error
		listOptions types.ListOptions
	)
	if err = httputils.ShouldBindAny(c, nil, nil, &listOptions); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	if r.Result, err = cr.c.Cluster().List(c, &listOptions); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	httputils.SetSuccess(c, r)
}
