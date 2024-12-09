package controller

import (
	"kubevulpes/cmd/app/config"
	"kubevulpes/pkg/db"
)

type VuplesInterface interface {
}

type vuples struct {
	cc      config.Config
	factory db.ShareDaoFactory
}

func New(cfg config.Config, f db.ShareDaoFactory) VuplesInterface {
	return &vuples{
		cc:      cfg,
		factory: f,
	}
}
