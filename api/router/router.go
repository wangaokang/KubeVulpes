package router

import (
	option "kubevulpes/cmd/app/options"
)

type RegisterFunc func(o *option.Options)

func InstallRouters(o *option.Options) {
	fs := []RegisterFunc{}

	install(o, fs...)
}

func install(o *option.Options, fs ...RegisterFunc) {
	for _, f := range fs {
		f(o)
	}
}
