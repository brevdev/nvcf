package function

import (
	"fmt"
	"strings"

	"github.com/brevdev/nvcf/config"
	"github.com/spf13/cobra"
)

// FunctionCmd returns a cobra.Command for managing NVIDIA Cloud Functions.
// It provides subcommands for various operations such as listing, creating,
// retrieving, and deleting functions, as well as running smoke tests.
//
// The command includes persistent pre-run checks to ensure the user is
// authenticated before executing any subcommands (except for the 'auth' command).
func FunctionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "function",
		Aliases: []string{"functions", "fn", "fns", "cf"},
		Short:   "Manage NVIDIA Cloud Functions",
		Long: `Create, list, update, call, deploy, and delete NVIDIA Cloud Functions. 
This command provides a comprehensive interface for interacting with 
NVIDIA Cloud Functions, allowing you to perform various operations 
on your serverless functions.`,
	}
	// Add subcommands
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

func authCheck(cmd *cobra.Command, args []string) error {
	config.Init()
	if !config.IsAuthenticated() {
		return fmt.Errorf("you are not authenticated. Please run 'nvcf auth login' first")
	}
	return nil
}

func shouldApplyAuthCheck(cmd *cobra.Command) bool {
	return cmd.Name() != "smoketest" && !strings.HasPrefix(cmd.Name(), "auth")
}
