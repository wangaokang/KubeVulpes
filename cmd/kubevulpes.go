package main

import (
	"io"
	"k8s.io/klog/v2"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/util/rand"

	"kubevulpes/cmd/app"
)

var (
	version = "v0.0.1-dev"
)

func main() {
	klog.InitFlags(nil)
	rand.Seed(time.Now().UnixNano())

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	cmd := app.NewServerCommand(version)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}

}
