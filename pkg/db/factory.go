package db

import (
	"KubeVulpes/pkg/db/user"
	"gorm.io/gorm"
)

type ShareDaoFactory interface {
	User() user.UserInterface
}

type shareDaoFactory struct {
	db *gorm.DB
}

func (f *shareDaoFactory) User() user.UserInterface { return user.NewUser(f.db) }

// 初始化一个调用数据库的接口
func NewDaoFactory(db *gorm.DB) ShareDaoFactory {
	return &shareDaoFactory{
		db: db,
	}
}
