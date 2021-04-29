package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "version",
		Aliases: []string{"v"},
		Short: "Print the version number of CDF upgrade",
		Long:  `All software has versions. This is CDF upgrade's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("CDF upgrade Static Site Generator v0.9 -- HEAD")
		},
	}

	return cmd
}
