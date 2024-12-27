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

package option

import (
	"fmt"
	"os"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/elastic/go-ucfg"
	"github.com/elastic/go-ucfg/yaml"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"kubevulpes/cmd/app/config"
	"kubevulpes/pkg/controller"
	"kubevulpes/pkg/db"
	vulpesModel "kubevulpes/pkg/db/model"
)

const (
	maxIdleConns      = 10
	maxOpenConns      = 100
	defaultListen     = 8080
	defaultConfigFile = "/etc/kubevulpes/config.yaml"
	defaultTokenKey   = "vuples"

	defaultSlowSQLDuration = 1 * time.Second

	rulesTableName = "rules"
)

type Options struct {
	// The default values.
	ComponentConfig config.Config
	HttpEngine      *gin.Engine

	// 数据库接口
	db      *gorm.DB
	Factory db.ShareDaoFactory

	// 控制器接口
	Controller controller.VuplesInterface

	//configFile 文件
	ConfigFile string

	// Authorization enforcement and policy management
	Enforcer *casbin.SyncedEnforcer
}

func NewOptions() (*Options, error) {
	return &Options{
		HttpEngine: gin.Default(), // 初始化默认 api 路由
		ConfigFile: defaultConfigFile,
	}, nil
}

func (o *Options) Complete() error {
	// todo 获取配置文件参数
	// 配置文件优先级: 默认配置，环境变量，命令行
	if len(o.ConfigFile) == 0 {
		// Try to read config file path from env.
		if cfgFile := os.Getenv("ConfigFile"); cfgFile != "" {
			o.ConfigFile = cfgFile
		} else {
			o.ConfigFile = defaultConfigFile
		}
	}

	// 解析配置文件
	if err := o.Binding(o.ConfigFile, &o.ComponentConfig); err != nil {
		return err
	}

	if o.ComponentConfig.Default.Listen == 0 {
		o.ComponentConfig.Default.Listen = defaultListen
	}
	if len(o.ComponentConfig.Default.JWTKey) == 0 {
		o.ComponentConfig.Default.JWTKey = defaultTokenKey
	}

	// 注册依赖组件
	if err := o.register(); err != nil {
		return err
	}

	o.Controller = controller.New(o.ComponentConfig, o.Factory, o.Enforcer)
	return nil
}

func (o *Options) Binding(configFile string, conf *config.Config) error {
	configContent, err := yaml.NewConfigWithFile(configFile, ucfg.PathSep("."))
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("config file not found")
		}
		return err
	}
	// 使用此解析参数不能使用"(Mode   Mode `config:"mode"`)" 自定义类型作为 config值，否者会导致此处参数无法解析，直接终止程序运行
	if err = configContent.Unpack(&conf); err != nil {
		return err
	}

	return nil
}

// BindFlags binds the kubevulpes Configuration struct fields
func (o *Options) BindFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.ConfigFile, "configfile", "", "The location of the kubevulpes configuration file")
}

func (o *Options) register() error {
	// 注册数据库
	if err := o.registerDatabase(); err != nil {
		return err
	}

	// TODO: 注册其他依赖
	if err := o.registerEnforcer(); err != nil {
		return err
	}
	// 初始化数据库日志
	o.ComponentConfig.Default.LogOptions.Init()

	return nil
}

func (o *Options) registerDatabase() error {
	sqlConfig := o.ComponentConfig.DB
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		sqlConfig.User,
		sqlConfig.Password,
		sqlConfig.Host,
		sqlConfig.Port,
		sqlConfig.Name)

	opt := &gorm.Config{
		Logger: db.NewLogger(logger.Info, defaultSlowSQLDuration),
	}

	dbCon, err := gorm.Open(mysql.Open(dsn), opt)
	if err != nil {
		return err
	}

	// 保存数据库对象
	o.db = dbCon

	// 设置数据库连接池
	sqlDB, err := dbCon.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)

	o.Factory, err = db.NewDaoFactory(dbCon, o.ComponentConfig.Default.AutoMigrate)
	if err != nil {
		return err
	}
	return nil
}

// This panics if o.db is nil.
func (o *Options) registerEnforcer() error {
	// Casbin 注册 casbin 权限表 注意 o.db 是否是 nil
	if !o.ComponentConfig.Default.AutoMigrate {
		return nil
	}
	a, err := gormadapter.NewAdapterByDBUseTableName(o.db, "", rulesTableName)
	if err != nil {
		return err
	}

	m, err := model.NewModelFromString(vulpesModel.RBACModel)
	if err != nil {
		return err
	}

	if o.Enforcer, err = casbin.NewSyncedEnforcer(m, a); err != nil {
		return err
	}

	// Add rbac policy.
	_, err = o.Enforcer.AddPolicy(vulpesModel.AdminPolicy.Raw())
	_, err = o.Enforcer.AddPolicy(vulpesModel.ReadOnlyPolicy.Raw())
	_, err = o.Enforcer.AddPolicy(vulpesModel.ReadWritePolicy.Raw())
	_, err = o.Enforcer.AddPolicy(vulpesModel.ReadWriteUpdatePolicy.Raw())

	// Add CustomKeyMatch function
	o.Enforcer.AddFunction("keyMatch2", vulpesModel.CustomKeyMatch)

	return err
}
