package db

import "gorm.io/gorm"

type ShareDaoFactory interface {
	User() UserInterface
}

type shareDaoFactory struct {
	db *gorm.DB
}

func (f *shareDaoFactory) User() UserInterface { return newUser(f.db) }

func NewDaoFactory(db *gorm.DB, migrate bool) (ShareDaoFactory, error) {
	if migrate {
		// 自动创建指定模型的数据库表结构
		if err := newMigrator(db).AutoMigrate(); err != nil {
			return nil, err
		}
	}

	return &shareDaoFactory{
		db: db,
	}, nil
}
