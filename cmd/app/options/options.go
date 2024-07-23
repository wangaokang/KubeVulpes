package options

import (
	"github.com/gin-gonic/gin"
	"kubevulpes/cmd/app/config"
)

type Options struct {
	// The default values.
	ComponentConfig config.Config
	HttpEngine      *gin.Engine

	//todo 预留数据库接口

	//configFile 文件
	ConfigFile string
}
