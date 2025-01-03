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

package router

import (
	"kubevulpes/api/middleware"
	"kubevulpes/api/router/audit"
	"kubevulpes/api/router/auth"
	"kubevulpes/api/router/cluster"
	"kubevulpes/api/router/user"
	option "kubevulpes/cmd/app/options"
)

type RegisterFunc func(o *option.Options)

func InstallRouters(o *option.Options) {
	fs := []RegisterFunc{
		middleware.InstallMiddlewares,
		cluster.NewRouter,
		user.NewRouter,
		audit.NewRouter,
		auth.NewRouter, // TODO: add auth router
	}

	install(o, fs...)
}

func install(o *option.Options, fs ...RegisterFunc) {
	for _, f := range fs {
		f(o)
	}
}
