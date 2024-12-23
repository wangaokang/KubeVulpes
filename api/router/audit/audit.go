package audit

import (
	"github.com/gin-gonic/gin"

	option "kubevulpes/cmd/app/options"
	"kubevulpes/pkg/controller"
)

type auditRouter struct {
	c controller.VuplesInterface
}

func NewRouter(o *option.Options) {
	r := auditRouter{c: o.Controller}
	r.initRouter(o.HttpEngine)
}

func (r *auditRouter) initRouter(httpEngine *gin.Engine) {
	auditRoute := httpEngine.Group("/api/v1/audits")
	{
		auditRoute.GET("", r.listAudit)
		auditRoute.GET("/:auditId", r.getAudit)
	}
}
