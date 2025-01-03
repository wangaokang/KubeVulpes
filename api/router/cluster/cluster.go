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

	option "kubevulpes/cmd/app/options"
	"kubevulpes/pkg/controller"
)

type clusterRouter struct {
	c controller.VuplesInterface
}

func NewRouter(o *option.Options) {
	r := clusterRouter{c: o.Controller}
	r.initRouter(o.HttpEngine)
}

func (r *clusterRouter) initRouter(httpEngine *gin.Engine) {
	clusterRoute := httpEngine.Group("/api/vulpes/clusters")
	{
		clusterRoute.POST("", r.createCluster)
		clusterRoute.GET("", r.listCluster)
		clusterRoute.GET("/:clusterId", r.getCluster)
		clusterRoute.DELETE("", r.deleteCluster)
		clusterRoute.PUT("/:clusterId", r.updateCluster)
	}
}
