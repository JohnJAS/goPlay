package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

func newInfraCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "infra",
		Aliases: []string{"i"},
		Short:   "Upgrade CDF infrastructure",
		Long: `
***********************************************************************************
   WARNING: This step is used to upgrade CDF infrastructure to 2021.05 release.
            The upgrade process is irreversible. You can NOT roll back.
            Make sure that all nodes in your cluster are in Ready status.
            Make sure that all Pods and Services are Running.

   NOTE:    You need to provide one temporary directory for saving the files for
            upgrade.
            Please make sure the directory has sufficient free space.

***********************************************************************************
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Pretend to upgrade CDF infrastructure.")
			return nil
		},
	}

	var drain bool
	var drainTimeOut string

	cmd.Flags().BoolVar(&drain, "drain", false, "Drain node before upgrade. It only takes effect during executing upgrade infra on worker nodes.")
	cmd.Flags().StringVar(&drainTimeOut, "drain-timeout", "3600", "The length of time(second) to wait before giving up to drain the node.")

	return cmd
}
