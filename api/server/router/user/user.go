package user

import "github.com/gin-gonic/gin"

type userRouter struct{}

func NewRouter(ginEngine *gin.Engine) {
	u := &userRouter{}
	u.initRoutes(ginEngine)
}

func (u *userRouter) initRoutes(ginEngine *gin.Engine) {
	userRoute := ginEngine.Group("/users")
	{
		userRoute.POST("", u.createUser)
		userRoute.DELETE("/:id", u.deleteUser)
		userRoute.PUT("/:id", u.updateUser)
		userRoute.GET("/:id", u.getUser)
		userRoute.GET("", u.listUsers)

	}
}
