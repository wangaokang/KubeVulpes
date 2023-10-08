package user

import (
	KubeVulpesmeta "KubeVulpes/api/meta"
	"KubeVulpes/pkg/db/model"
	dberrors "KubeVulpes/pkg/errors"
	"context"
	"gorm.io/gorm"
	"time"
)

type UserInterface interface {
	Create(ctx context.Context, obj *model.User) (*model.User, error)
	Update(ctx context.Context, uid int64, updates map[string]interface{}) error
	Delete(ctx context.Context, uid int64) error
	Get(ctx context.Context, uid int64) (*model.User, error)
	List(ctx context.Context, selector *KubeVulpesmeta.ListSelector) ([]*model.User, error)
	// TODO: Add more methods
}

type user struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) UserInterface {
	return &user{db}
}

func (u *user) Create(ctx context.Context, obj *model.User) (*model.User, error) {
	now := time.Now()
	obj.CreateAt = now
	obj.UpdateAt = now
	if err := u.db.Create(obj).Error; err != nil {
		return nil, err
	}

	return obj, nil
}

func (u *user) Update(ctx context.Context, uid int64, updates map[string]interface{}) error {
	// 系统维护字段
	updates["updateAt"] = time.Now()

	f := u.db.Model(&model.User{}).
		Where("id = ? ", uid).
		Updates(updates)
	if f.Error != nil {
		return f.Error
	}

	if f.RowsAffected == 0 {
		return dberrors.ErrRecordNotUpdate
	}

	return nil
}

func (u *user) Delete(ctx context.Context, uid int64) error {
	return u.db.
		Where("id = ?", uid).
		Delete(&model.User{}).
		Error
}

func (u *user) Get(ctx context.Context, uid int64) (*model.User, error) {
	var obj model.User
	if err := u.db.Where("id = ?", uid).First(&obj).Error; err != nil {
		return nil, err
	}

	return &obj, nil
}

// List 分页查询
func (u *user) List(ctx context.Context, selector *KubeVulpesmeta.ListSelector) ([]*model.User, error) {
	var (
		users []*model.User
	)
	if err := u.db.Limit(selector.Limit).Offset((selector.Page - 1) * selector.Limit).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
