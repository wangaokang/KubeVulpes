package main

import (
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/klog/v2"

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
