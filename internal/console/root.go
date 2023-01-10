package console

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "cobra-example",
	Short: "An example of cobra",
	Long: `This application shows how to create modern CLI
			applications in go using Cobra CLI library`,
}

// Execute runs the root command for the application and handles any errors that may occur during execution.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
