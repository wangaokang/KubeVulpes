package auth

import (
	"github.com/gin-gonic/gin"

	"kubevulpes/api/httputils"
	"kubevulpes/pkg/types"
)

func (a *authRouter) listPolicy(c *gin.Context) {
	r := httputils.NewResponse()
	var (
		req types.ListRBACPolicyRequest
		err error
	)
	if err = c.ShouldBindQuery(&req); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	//policies, err := a.c.Auth().ListRBACPolicy(c)
	//if err != nil {
	//	httputils.SetFailed(c, r, err)
	//	return
	//}
	//r.Data = policies
	httputils.SetSuccess(c, r)
}

func (a *authRouter) createPolicy(c *gin.Context) {
	r := httputils.NewResponse()
	var req types.RBACPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	//if err := a.c.Auth().CreateRBACPolicy(c, &req); err != nil {
	//	httputils.SetFailed(c, r, err)
	//	return
	//}

	httputils.SetSuccess(c, r)
}

func (a *authRouter) deletePolicy(c *gin.Context) {
	r := httputils.NewResponse()
	var req types.RBACPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	//if err := a.c.Auth().DeleteRBACPolicy(c); err != nil {
	//	httputils.SetFailed(c, r, err)
	//	return
	//}
	httputils.SetSuccess(c, r)
}

func (a *authRouter) listGroupBinding(c *gin.Context) {
	r := httputils.NewResponse()
	var (
		req types.ListGroupBindingRequest
		err error
	)
	if err = c.ShouldBindQuery(&req); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	//if r.Result, err = a.c.Auth().ListGroupBindings(c, &req); err != nil {
	//	httputils.SetFailed(c, r, err)
	//	return
	//}

	httputils.SetSuccess(c, r)
}

func (a *authRouter) createGroupBinding(c *gin.Context) {
	r := httputils.NewResponse()
	var req types.GroupBindingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	//if err := a.c.Auth().CreateGroupBinding(c, &req); err != nil {
	//	httputils.SetFailed(c, r, err)
	//	return
	//}

	httputils.SetSuccess(c, r)
}

func (a *authRouter) deleteGroupBinding(c *gin.Context) {
	r := httputils.NewResponse()
	var req types.GroupBindingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	//if err := a.c.Auth().DeleteGroupBinding(c, &req); err != nil {
	//	httputils.SetFailed(c, r, err)
	//	return
	//}

	httputils.SetSuccess(c, r)
}
