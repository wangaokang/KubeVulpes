package user

import (
	"KubeVulpes/api/server/httputils"
	"KubeVulpes/api/types"
	"KubeVulpes/pkg/kubeVulpes"
	"github.com/gin-gonic/gin"
)

func (u *userRouter) createUser(c *gin.Context) {
	r := httputils.NewResponse()
	var user types.User
	if err := c.ShouldBindJSON(&user); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	if err := kubeVulpes.CoreV1.User().Create(c, &user); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
}

func (u *userRouter) deleteUser(c *gin.Context) {

}

func (u *userRouter) updateUser(c *gin.Context) {

}

func (u *userRouter) getUser(c *gin.Context) {

}

func (u *userRouter) listUsers(c *gin.Context) {

}
