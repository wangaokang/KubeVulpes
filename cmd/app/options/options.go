package options

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"kubevulpes/cmd/app/config"
)

const (
	maxIdleConns      = 10
	maxOpenConns      = 100
	defaultConfigFile = "/etc/kubevulpes/config.yaml"
)

type Options struct {
	// The default values.
	ComponentConfig config.Config
	HttpEngine      *gin.Engine

	//todo 预留数据库接口

	//configFile 文件
	ConfigFile string
}

func NewOptions() (*Options, error) {
	return &Options{
		HttpEngine: gin.Default(), // 初始化默认 api 路由
		ConfigFile: defaultConfigFile,
	}, nil
}

func (o *Options) Complete() error {

	return nil
}

// BindFlags binds the pixiu Configuration struct fields
func (o *Options) BindFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.ConfigFile, "configfile", "", "The location of the kubevulpes configuration file")
}

func (o *Options) register() error {
	// 注册数据库
	if err := o.registerDatabase(); err != nil {
		return err
	}

	// TODO: 注册其他依赖
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

	opt := &gorm.Config{}
	if o.ComponentConfig.Default.Mode == "debug" {
		opt.Logger = logger.Default.LogMode(logger.Info)
	}

	DB, err := gorm.Open(mysql.Open(dsn), opt)
	if err != nil {
		return err
	}
	// 设置数据库连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)

	//o.Factory, err = db.NewDaoFactory(DB, o.ComponentConfig.Default.AutoMigrate)
	if err != nil {
		return err
	}
	return nil
}
