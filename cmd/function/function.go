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
		Aliases: []string{"functions", "fn", "fns", "cf"},
		Short:   "Manage NVIDIA Cloud Functions",
		Long: `Create, list, update, call, deploy, and delete NVIDIA Cloud Functions. 
This command provides a comprehensive interface for interacting with 
NVIDIA Cloud Functions, allowing you to perform various operations 
on your serverless functions.`,
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
	cmd.AddCommand(functionStopCmd())
	cmd.AddCommand(functionWatchCmd())

	return cmd
}
