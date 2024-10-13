package auth

import (
	"fmt"
	"os"

	"github.com/brevdev/nvcf/config"
	"github.com/spf13/cobra"
)

func AuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication for the CLI",
		Long:  `Authenticate with NVIDIA Cloud and configure the CLI to use your API key.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			config.Init()
			if cmd.Parent().Name() != "auth" && !config.IsAuthenticated() {
				fmt.Println("You are not authenticated. Please run 'nvcf auth login' first.")
				os.Exit(1)
			}
		},
	}

	cmd.AddCommand(authLoginCmd())
	cmd.AddCommand(authLogoutCmd())
	cmd.AddCommand(authStatusCmd())
	cmd.AddCommand(authConfigureDockerCmd())

	cmd.AddCommand(authWhoAmICmd())
	cmd.AddCommand(authOrgsCmd())
	cmd.AddCommand(authOrgIDCmd())

	return cmd
}
