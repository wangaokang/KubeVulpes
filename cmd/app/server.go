package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"kubevulpes/api/router"
	option "kubevulpes/cmd/app/options"
)

func NewServerCommand(version string) *cobra.Command {
	opts, err := option.NewOptions()
	if err != nil {
		log.Fatalf("unable to initialize command options: %v", err)
	}
	cmd := &cobra.Command{
		Use:  "kube-vulpes",
		Long: `kube-vulpes build`,
		Run: func(cmd *cobra.Command, args []string) {
			if err = opts.Complete(); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
			if err = Run(opts); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}

	// 绑定命令行参数
	opts.BindFlags(cmd)

	verCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Long:  "Print version and exit.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}
	cmd.AddCommand(verCmd)
	return cmd
}

// 优雅启动服务
func Run(opt *option.Options) error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", opt.ComponentConfig.Default.Listen),
		Handler: opt.HttpEngine,
	}

	// 安装 http 路由
	router.InstallRouters(opt)

	// Initializing the server in a goroutine so that it won't block the graceful shutdown handling below
	go func() {
		klog.Info("starting license server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			klog.Fatal("failed to listen license server: ", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting vuples server down ...")

	// The context is used to inform the server it has 5 seconds to finish the request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("vuples server forced to shutdown: %v", err)
	}
	return nil
}
