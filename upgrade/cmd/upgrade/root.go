package main

import (
	"github.com/spf13/cobra"
)

func newRootCmd() (*cobra.Command, error) {
	// rootCmd represents the base command when called without any subcommands
	cmd := &cobra.Command{
		Use:   "upgrade.sh  [-i|--infra ] | [-u|--upgrade] | [-c|--clean] [Options]",
		Short: "Upgrade CDF infrastructure and components",
		Long:  `./upgrade.sh  [-i|--infra ] | [-u|--upgrade] | [-c|--clean] [Options]`,
	}

	cmd.AddCommand(newInfraCmd())
	cmd.AddCommand(newVersionCmd())

	return cmd, nil
}
