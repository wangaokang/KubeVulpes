/*
Copyright 2024 The Vuples Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package model

import "kubevulpes/pkg/db/model/base"

type UserRole uint8

const (
	RoleReader          UserRole = iota * 2 // 只读用户
	RoleReadWrite                           // 读写用户
	RoleReadWriteUpdate                     // 读写改用户
	RoleAdmin                               // 管理员 只能管理自己创建的用户
	RoleRoot                                // 超级管理员 可以管理所有用户
)

const (
	UserStatusNormalUserStatus UserStatus = iota // 正常
	UserStatusDisabled                           // 禁用
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
	IsDeleted   bool       `gorm:"column:is_deleted;default:false" json:"is_deleted"`
}

func init() {
	register(&User{})
}

func (user *User) TableName() string {
	return "users"
}
