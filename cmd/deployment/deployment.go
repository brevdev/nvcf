package deployment

import (
	"fmt"
	"os"

	"github.com/brevdev/nvcf/config"
	"github.com/spf13/cobra"
)

func DeploymentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deployment",
		Aliases: []string{"deploy", "d"},
		Short:   "Manage NVIDIA Cloud Deployments",
		Long:    `Create, list, update, call, deploy, and delete NVIDIA Cloud Deployments`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			config.Init()
			if cmd.Name() != "auth" && !config.IsAuthenticated() {
				fmt.Println("You are not authenticated. Please run 'nvcf auth login' first.")
				os.Exit(1)
			}
		},
	}

	cmd.AddCommand(deploymentListCmd())
	cmd.AddCommand(deploymentGetCmd())
	return cmd
}
