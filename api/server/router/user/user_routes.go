package user

import (
	KubeVulpesmeta "KubeVulpes/api/meta"
	"KubeVulpes/api/server/httputils"
	"KubeVulpes/api/types"
	"KubeVulpes/pkg/kubeVulpes"
	"KubeVulpes/pkg/utils"
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

	httputils.SetSuccess(c, r)
}

func (u *userRouter) deleteUser(c *gin.Context) {
	r := httputils.NewResponse()
	uid, err := utils.ParseInt64(c.Param("id"))
	if err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	if err := kubeVulpes.CoreV1.User().Delete(c, uid); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	httputils.SetSuccess(c, r)
}

func (u *userRouter) updateUser(c *gin.Context) {
	r := httputils.NewResponse()
	var (
		err  error
		user types.User
	)
	if err = c.ShouldBindJSON(&user); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	user.Id, err = utils.ParseInt64(c.Param("id"))
	if err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	if err := kubeVulpes.CoreV1.User().Update(c, &user); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	httputils.SetSuccess(c, r)
}

func (u *userRouter) getUser(c *gin.Context) {
	r := httputils.NewResponse()
	uid, err := utils.ParseInt64(c.Param("id"))
	if err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	r.Result, err = kubeVulpes.CoreV1.User().Get(c, uid)
	if err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	httputils.SetSuccess(c, r)
}

func (u *userRouter) listUsers(c *gin.Context) {
	r := httputils.NewResponse()
	var err error
	if r.Result, err = kubeVulpes.CoreV1.User().List(c, KubeVulpesmeta.ParseListSelector(c)); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	httputils.SetSuccess(c, r)
}
