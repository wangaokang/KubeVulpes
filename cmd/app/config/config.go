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

package config

import (
	logutil "kubevulpes/pkg/util/log"
)

type Config struct {
	DB      DBOptions      `config:"db"`
	Default DefaultOptions `config:"default"`
}

type DBOptions struct {
	Name     string `config:"name"`
	Host     string `config:"host"`
	User     string `config:"user"`
	Password string `config:"password"`
	Port     int    `config:"port"`
}

type DefaultOptions struct {
	Listen int    `config:"listen"`
	Mode   string `config:"mode"` // 此处使用 unpack 渲染不能自定义变量
	JWTKey string `config:"jwt_key"`

	// 自动创建指定模型的数据库表结构，不会更新已存在的数据库表
	AutoMigrate bool `config:"auto_migrate"`

	logutil.LogOptions `config:",inline"`
}

func (d *DefaultOptions) InDebug() bool {
	return d.Mode == "debug"
}
