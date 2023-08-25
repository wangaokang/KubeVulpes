package option

import (
	"KubeVulpes/cmd/app/config"
	"KubeVulpes/pkg/db"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	maxIdleConns = 10
	maxOpenConns = 100

	defaultConfigFile = "/etc/KubeVulpus/config.yaml"
)

type Options struct {
	// The default values.
	ComponentConfig config.Config
	GinEngine       *gin.Engine
	DB              *gorm.DB           //数据库设置操作
	Factory         db.ShareDaoFactory // 数据库接口

	// ConfigFile is the location of the KubeVulpus server's configuration file.
	ConfigFile string
}

type Config struct {
	configFile string
	configType string

	data []byte
}

func (o *Options) registerDatabase() error {
	sqlConfig := o.ComponentConfig.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		sqlConfig.User,
		sqlConfig.Password,
		sqlConfig.Host,
		sqlConfig.Port,
		sqlConfig.Name)

	var err error
	if o.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		return err
	}

	return nil
}
