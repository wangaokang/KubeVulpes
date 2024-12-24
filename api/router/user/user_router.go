package user

import (
	"github.com/gin-gonic/gin"

	"kubevulpes/api/httputils"
	"kubevulpes/pkg/types"
)

type IdMeta struct {
	UserId int64 `uri:"userId" binding:"required"`
}

func (u *userRouter) createUser(c *gin.Context) {
	r := httputils.NewResponse()

	var (
		req types.CreateUserRequest
		err error
	)
	if err = c.ShouldBindJSON(&req); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	//if err = u.c.User().Create(c, &req); err != nil {
	//	httputils.SetFailed(c, r, err)
	//	return
	//}

	httputils.SetSuccess(c, r)
}

func (u *userRouter) deleteUser(c *gin.Context) {
	r := httputils.NewResponse()
	var (
		idMeta IdMeta
		err    error
	)
	if err = httputils.ShouldBindAny(c, nil, &idMeta, nil); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	//if err = u.c.User().Delete(c, idMeta.UserId); err != nil {
	//	httputils.SetFailed(c, r, err)
	//	return
	//}

	httputils.SetSuccess(c, r)
}

func (u *userRouter) getUser(c *gin.Context) {
	r := httputils.NewResponse()
	// todo
	httputils.SetSuccess(c, r)
}

func (u *userRouter) listUser(c *gin.Context) {
	r := httputils.NewResponse()
	// todo
	httputils.SetSuccess(c, r)
}

func (u *userRouter) updateUser(c *gin.Context) {
	r := httputils.NewResponse()
	// todo
	httputils.SetSuccess(c, r)
}
