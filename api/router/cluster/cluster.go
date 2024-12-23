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
	clusterRoute := httpEngine.Group("/api/v1/clusters")
	{
		clusterRoute.POST("", r.createCluster)
		clusterRoute.GET("", r.listCluster)
		clusterRoute.GET("/:clusterId", r.getCluster)
		clusterRoute.DELETE("", r.deleteCluster)
		clusterRoute.PUT("/:clusterId", r.updateCluster)
	}
}
