/*
Copyright 2024 The Vuples Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
	if err = u.c.User().Create(c, &req); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

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
