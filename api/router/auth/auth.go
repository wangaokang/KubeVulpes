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

package auth

import (
	"github.com/gin-gonic/gin"

	option "kubevulpes/cmd/app/options"
	"kubevulpes/pkg/controller"
)

const (
	AuthBashPath   = "/api/vulpes/auth"
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
