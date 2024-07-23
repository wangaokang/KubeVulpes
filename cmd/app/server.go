package app

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewServerCommand() *cobra.Command {
	//option := options.NewServerRunOption()
	cmd := &cobra.Command{
		Use:  "kubeV",
		Long: "kubevupels is a ops tools for kubernetes.",
		Run:  func(cmd *cobra.Command, args []string) {},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}
}
