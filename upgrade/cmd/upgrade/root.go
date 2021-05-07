package upgrade

import (
	"github.com/spf13/cobra"
	"io"
)

func NewRootCmd(io io.Writer) (*cobra.Command, error) {
	// rootCmd represents the base command when called without any subcommands
	cmd := &cobra.Command{
		Use:     "upgrade",
		Example: "upgrade {i|infra} | {u|upgrade} | {c|clean} [Flags]",
		Short:   "CDF upgrade",
		Long:    `Upgrade CDF infrastructure and components`,
	}

	//add subcommand
	cmd.AddCommand(newInfraCmd())
	cmd.AddCommand(newComponentCmd())
	cmd.AddCommand(newCleanCommand())
	cmd.AddCommand(newVersionCmd(io))
	cmd.SetHelpCommand(&cobra.Command{Hidden: true})

	//add global flags
	var confirm bool
	var backupFolder string

	cmd.PersistentFlags().BoolVarP(&confirm, "yes", "y", false, "Answer yes for any confirmations")
	cmd.PersistentFlags().StringVarP(&backupFolder, "temp", "t", "/tmp", "Specify an absolute path that already exists for storing CDF upgrade backup files")

	return cmd, nil
}
