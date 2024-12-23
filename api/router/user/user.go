package user

import (
	"github.com/gin-gonic/gin"

	option "kubevulpes/cmd/app/options"
	"kubevulpes/pkg/controller"
)

type userRouter struct {
	c controller.VuplesInterface
}

func NewRouter(o *option.Options) {
	r := &userRouter{c: o.Controller}
	r.initRouter(o.HttpEngine)
}

func (u *userRouter) initRouter(httpEngine *gin.Engine) {
	userRoute := httpEngine.Group("/api/v1/users")
	{
		userRoute.POST("", u.createUser)
		userRoute.GET("", u.listUser)
		userRoute.GET("/:userId", u.getUser)
		userRoute.DELETE("", u.deleteUser)
		userRoute.PUT("", u.updateUser)
	}
}
