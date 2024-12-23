package router

import (
	"kubevulpes/api/router/audit"
	"kubevulpes/api/router/cluster"
	"kubevulpes/api/router/role"
	"kubevulpes/api/router/user"
	option "kubevulpes/cmd/app/options"
)

type RegisterFunc func(o *option.Options)

func InstallRouters(o *option.Options) {
	fs := []RegisterFunc{
		cluster.NewRouter,
		role.NewRouter,
		user.NewRouter,
		audit.NewRouter,
		//auth.NewRouter,  // TODO: add auth router
	}

	install(o, fs...)
}

func install(o *option.Options, fs ...RegisterFunc) {
	for _, f := range fs {
		f(o)
	}
}
