package upgrade

import (
	"fmt"
	"github.com/spf13/cobra"
)

func newCleanCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "clean",
		Aliases: []string{"c"},
		Short:   "Clean useless runtime images after CDF upgrade",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Pretend to clean runtime images")
			return nil
		},
	}

	return cmd
}