package user

import (
	"KubeVulpes/pkg/db/model"
	"context"
	"gorm.io/gorm"
	"time"
)

type UserInterface interface {
	Create(ctx context.Context, obj *model.User) (*model.User, error)
	Update(ctx context.Context, obj *model.User, uid int) (*model.User, error)
	Delete(ctx context.Context, uid int) error
	Get(ctx context.Context, uid int) (*model.User, error)
	List(ctx context.Context, page int, pageSize int) ([]*model.User, int, error)
	// TODO: Add more methods
}

type user struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) UserInterface {
	return &user{db}
}

func (u *user) Create(ctx context.Context, obj *model.User) (*model.User, error) {
	obj.CreateAt = time.Now()
	if err := u.db.Create(obj).Error; err != nil {
		return nil, err
	}

	return obj, nil
}

func (u *user) Update(ctx context.Context, obj *model.User, uid int) (*model.User, error) {
	obj.UpdateAt = time.Now()
	f := u.db.Model(&model.User{}).Where("id = ?", uid).Updates(obj)
	if f.Error != nil {
		return nil, f.Error
	}

	return obj, nil
}

func (u *user) Delete(ctx context.Context, uid int) error {
	return u.db.
		Where("id = ?", uid).
		Delete(&model.User{}).
		Error
}

func (u *user) Get(ctx context.Context, uid int) (*model.User, error) {
	var obj model.User
	if err := u.db.Where("id = ?", uid).First(&obj).Error; err != nil {
		return nil, err
	}

	return &obj, nil
}

// List 分页查询
func (u *user) List(ctx context.Context, page int, pageSize int) ([]*model.User, int, error) {
	var (
		users []*model.User
		total int64
	)
	if err := u.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := u.db.Limit(pageSize).Offset((page - 1) * pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, 0, nil
}
