package debug

import (
	"fmt"

	"github.com/brevdev/nvcf/cmd/preflight/brev"
	"github.com/spf13/cobra"
)

func debugStopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop a debug environment",
		Long:  `Stop and clean up a debug environment for an NVCF function`,
		RunE:  runDebugStop,
	}

	cmd.Flags().StringP("function-id", "f", "", "The ID of the function to debug")

	return cmd
}

func runDebugStop(cmd *cobra.Command, args []string) error {
	brevClient := brev.NewBrevClient()

	if !brevClient.IsBrevCLIInstalled() {
		return fmt.Errorf("brev CLI is not installed. Please install it first")
	}

	loggedIn, err := brevClient.IsLoggedIn()
	if err != nil {
		return fmt.Errorf("error checking if logged in: %w", err)
	}
	if !loggedIn {
		err := brevClient.Login()
		if err != nil {
			return fmt.Errorf("error logging in: %w", err)
		}
		fmt.Println("Logged in to Brev")
	}

	functionId, err := cmd.Flags().GetString("function-id")
	if err != nil {
		return fmt.Errorf("error getting function-id flag: %w", err)
	}

	if functionId == "" {
		return fmt.Errorf("function-id is required. Please provide a function-id using the -f or --function-id flag")
	}

	// Delete the debug instance
	err = brevClient.DeleteInstance(functionId)
	if err != nil {
		return fmt.Errorf("error deleting Brev instance: %w", err)
	}

	fmt.Printf("Successfully stopped and removed debug environment for function %s\n", functionId)

	return nil
}
