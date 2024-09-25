package debug

import (
	"fmt"

	"github.com/brevdev/nvcf/cmd/preflight/brev"
	"github.com/spf13/cobra"
)

func debugStopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop <function-id>",
		Short: "Stop a debug environment",
		Long:  `Stop and clean up a debug environment for an NVCF function`,
		Args:  cobra.ExactArgs(1),
		RunE:  runDebugStop,
	}

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

	functionId := args[0]

	// Delete the debug instance
	err = brevClient.DeleteInstance(functionId)
	if err != nil {
		return fmt.Errorf("error deleting Brev instance: %w", err)
	}

	fmt.Printf("Successfully stopped and removed debug environment for function %s\n", functionId)

	return nil
}
