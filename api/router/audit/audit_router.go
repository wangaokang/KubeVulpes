package audit

import (
	"github.com/gin-gonic/gin"

	"kubevulpes/api/httputils"
)

func (a *auditRouter) getAudit(c *gin.Context) {
	r := httputils.NewResponse()
	// todo
	httputils.SetSuccess(c, r)
}

func (a *auditRouter) listAudit(c *gin.Context) {
	r := httputils.NewResponse()
	// todo
	httputils.SetSuccess(c, r)
}
