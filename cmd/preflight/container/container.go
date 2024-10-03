package container

import (
	"github.com/spf13/cobra"
)

func ContainerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Aliases: []string{"c", "container"},
		Use:     "container",
		Short:   "Push a container to nvcr.io to start using it with NVCF",
		Long:    `Push a container to nvcr.io to start using it with NVCF`,
	}

	cmd.AddCommand(PushCmd())

	return cmd
}
