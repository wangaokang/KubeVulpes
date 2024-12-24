package model

import "kubevulpes/pkg/db/model/base"

type UserRole uint8

const (
	RoleReader          UserRole = iota * 2 // 只读用户
	RoleReadWrite                           // 读增用户
	RoleReadWriteUpdate                     // 读增改用户
	RoleAdmin                               // 管理员
	RoleRoot                                // 超级管理员
)

const (
	UserStatusDisabled         UserStatus = iota // 禁用
	UserStatusNormalUserStatus                   // 正常
	UserStatusDeleted                            // 删除
)

type UserStatus uint8 // TODO

type User struct {
	base.Model
	Name        string     `gorm:"types:varchar(128);column:name;not null" json:"name"`
	Password    string     `gorm:"types:varchar(256);column:password;not null" json:"password"`
	Role        UserRole   `gorm:"types:tinyint;column:role;not null" json:"role"` // 0
	Email       string     `gorm:"types:varchar(256);column:email;not null" json:"email"`
	Status      UserStatus `gorm:"types:tinyint;column:status;not null" json:"status"`
	Description string     `gorm:"types:text;column:description;not null" json:"description"`
}

func init() {
	register(&User{})
}

func (user *User) TableName() string {
	return "users"
}
