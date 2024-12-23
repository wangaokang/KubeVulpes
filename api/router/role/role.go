package role

import (
	"github.com/gin-gonic/gin"

	option "kubevulpes/cmd/app/options"
	"kubevulpes/pkg/controller"
)

type roleRouter struct {
	c controller.VuplesInterface
}

func NewRouter(o *option.Options) {
	r := roleRouter{c: o.Controller}
	r.initRouter(o.HttpEngine)
}

func (r *roleRouter) initRouter(httpEngine *gin.Engine) {
	roleRoute := httpEngine.Group("/api/v1/roles")
	{
		roleRoute.POST("", r.createRole)
		roleRoute.GET("", r.listRole)
		roleRoute.GET("/:roleId", r.getRole)
		roleRoute.DELETE("", r.deleteRole)
		roleRoute.PUT("/:roleId", r.updateRole)
	}
}
