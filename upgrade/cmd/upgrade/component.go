package upgrade

import (
	"fmt"
	"github.com/spf13/cobra"
)

func newComponentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "upgrade",
		Aliases: []string{"u"},
		Short:   "Upgrade CDF components",
		Long: `
***********************************************************************************
   WARNING: This step is used to upgrade CDF components to ${TARGET_RELEASE_VERSION_WITH_DOT} release. 
            The upgrade process is irreversible. You can NOT roll back.
            Make sure that all nodes in your cluster are in Ready status.
            Make sure that all Pods and Services are Running.

***********************************************************************************"
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Pretend to upgrade CDF components.")
			return nil
		},
	}

	return cmd
}
