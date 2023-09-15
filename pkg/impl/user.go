package impl

import (
	"KubeVulpes/api/types"
	"KubeVulpes/pkg/db"
	"KubeVulpes/pkg/db/model"
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type UserGetter interface {
	User() UserInterface
}

type UserInterface interface {
	Create(ctx context.Context, obj *types.User) error
	preCreate(ctx context.Context, obj *types.User) error
	Delete(ctx context.Context, uid int) error
	Get(ctx context.Context, uid int) (*types.User, error)
}

type user struct {
	factory db.ShareDaoFactory
}

func newUser(c *KubeVulpes) UserInterface {
	return &user{c.factory}
}

// 创建前检查：
// 1. 用户名不能为空
// 2. 用户密码不能为空
// 3. 其他创建前检查
func (u *user) preCreate(ctx context.Context, obj *types.User) error {
	if len(obj.Name) == 0 || len(obj.Password) == 0 {
		return fmt.Errorf("user name or password could not be empty")
	}

	return nil
}

func (u *user) Create(ctx context.Context, obj *types.User) error {
	if err := u.preCreate(ctx, obj); err != nil {
		return err
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(obj.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if _, err = u.factory.User().Create(ctx, &model.User{
		Name:        obj.Name,
		Password:    string(encryptedPassword),
		Role:        obj.Role,
		Status:      obj.Status,
		Email:       obj.Email,
		Description: obj.Description,
	}); err != nil {
		return err
	}
	return nil
}

func (u *user) Delete(ctx context.Context, uid int) error {
	if err := u.factory.User().Delete(ctx, uid); err != nil {
		return err
	}

	return nil
}

func (u *user) Get(ctx context.Context, uid int) (*types.User, error) {
	user, err := u.factory.User().Get(ctx, uid)
	if err != nil {
		return nil, err
	}

	return model2Type(user), nil
}

func model2Type(u *model.User) *types.User {
	return &types.User{
		Id:          u.Id,
		Name:        u.Name,
		Status:      u.Status,
		Role:        u.Role,
		Email:       u.Email,
		Description: u.Description,
		TimeOption:  types.FormatTime(u.CreateAt, u.UpdateAt),
	}
}
