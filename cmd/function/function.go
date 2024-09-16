package function

import (
	"fmt"
	"os"

	"github.com/brevdev/nvcf/config"
	"github.com/spf13/cobra"
)

func FunctionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "function",
		Aliases: []string{"fn", "fns", "cf"},
		Short:   "Manage NVIDIA Cloud Functions",
		Long:    `Create, list, update, call, deploy, and delete NVIDIA Cloud Functions. You can also specify a YAML file to create multiple functions at once.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			config.Init()
			if cmd.Name() != "auth" && !config.IsAuthenticated() {
				fmt.Println("You are not authenticated. Please run 'nvcf auth login' first.")
				os.Exit(1)
			}
		},
	}

	cmd.AddCommand(functionListCmd())
	cmd.AddCommand(functionCreateCmd())
	cmd.AddCommand(functionGetCmd())
	cmd.AddCommand(functionDeleteCmd())
	cmd.AddCommand(functionUpdateCmd())
	cmd.AddCommand(functionDeployCmd())

	return cmd
}
