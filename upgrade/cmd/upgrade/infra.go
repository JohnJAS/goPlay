package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

func newInfraCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "infra",
		Short: "Upgrade CDF infrastructure",
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
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("CDF upgrade Static Site Generator v0.9 -- HEAD")
		},
	}

	return cmd
}
