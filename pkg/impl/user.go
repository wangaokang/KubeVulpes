package impl

import (
	KubeVulpesmeta "KubeVulpes/api/meta"
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
	preCreate(ctx context.Context, obj *types.User) error
	Create(ctx context.Context, obj *types.User) error
	Update(ctx context.Context, obj *types.User) error
	Delete(ctx context.Context, uid int64) error
	Get(ctx context.Context, uid int64) (*types.User, error)
	List(ctx context.Context, selector *KubeVulpesmeta.ListSelector) ([]*types.User, error)
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

func (u *user) Update(ctx context.Context, obj *types.User) error {
	oldUser, err := u.factory.User().Get(ctx, obj.Id)
	if err != nil {
		return err
	}
	updates := u.parseUserUpdates(oldUser, obj)
	if err := u.factory.User().Update(ctx, oldUser.Id, updates); err != nil {
		return err
	}

	return nil
}

func (u *user) Delete(ctx context.Context, uid int64) error {
	if err := u.factory.User().Delete(ctx, uid); err != nil {
		return err
	}

	return nil
}

func (u *user) Get(ctx context.Context, uid int64) (*types.User, error) {
	user, err := u.factory.User().Get(ctx, uid)
	if err != nil {
		return nil, err
	}

	return model2Type(user), nil
}

func (u *user) List(ctx context.Context, selector *KubeVulpesmeta.ListSelector) ([]*types.User, error) {
	users, err := u.factory.User().List(ctx, selector)
	if err != nil {
		return nil, err
	}

	var ret []*types.User
	for _, user := range users {
		ret = append(ret, model2Type(user))
	}

	return ret, nil
}

func (u *user) parseUserUpdates(oldObj *model.User, newObj *types.User) map[string]interface{} {
	updates := make(map[string]interface{})

	if oldObj.Status != newObj.Status { // 更新状态
		updates["status"] = newObj.Status
	}
	if oldObj.Role != newObj.Role { // 更新用户角色
		updates["role"] = newObj.Role
	}
	if oldObj.Email != newObj.Email { // 更新邮件
		updates["email"] = newObj.Email
	}
	if oldObj.Description != newObj.Description { // 更新描述
		updates["description"] = newObj.Description
	}
	if oldObj.Password != newObj.Password { // 更新密码
		updates["password"] = newObj.Password
	}
	if oldObj.Name != newObj.Name { // 更新用户名
		updates["name"] = newObj.Name
	}

	return updates
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
