package auth

import (
	"github.com/gin-gonic/gin"

	option "kubevulpes/cmd/app/options"
	"kubevulpes/pkg/controller"
)

const (
	AuthBashPath   = "/api/v1/auth"
	PolicySubPath  = "/policy"
	BindingSubPath = "/binding"
)

type authRouter struct {
	c controller.VuplesInterface
}

func NewRouter(o *option.Options) {
	router := &authRouter{
		c: o.Controller,
	}
	router.initRoutes(o.HttpEngine)
}

func (a *authRouter) initRoutes(httpEngine *gin.Engine) {
	authRoute := httpEngine.Group(AuthBashPath)
	{
		policyRoute := authRoute.Group(PolicySubPath)
		policyRoute.POST("", a.createPolicy)
		policyRoute.GET("", a.listPolicy)
		policyRoute.DELETE("", a.deletePolicy)
	}
	{
		bindingRoute := authRoute.Group(BindingSubPath)
		bindingRoute.POST("", a.createGroupBinding)
		bindingRoute.GET("", a.listGroupBinding)
		bindingRoute.DELETE("", a.deletePolicy)
	}
}
