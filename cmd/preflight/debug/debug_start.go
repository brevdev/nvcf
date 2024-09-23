package debug

import (
	"fmt"

	"github.com/brevdev/nvcf/api"
	"github.com/spf13/cobra"
)

func debugStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start a debug environment",
		Long:  `Create and start a debug environment for an NVCF function`,
		RunE:  runDebugStart,
	}
}

func runDebugStart(cmd *cobra.Command, args []string) error {
	brevClient := api.NewBrevClient()

	if !brevClient.IsBrevCLIInstalled() {
		return fmt.Errorf("brev CLI is not installed. Please install it first")
	}

	loggedIn, err := brevClient.IsLoggedIn()
	if err != nil {
		return err
	}

	if !loggedIn {
		fmt.Println("You are not logged in. Starting Brev login process...")
		err = brevClient.Login()
		if err != nil {
			return err
		}
		fmt.Println("Successfully logged in with Brev CLI")
	}
	// TODO: Implement debug environment setup logic here
	fmt.Println("Setting up debug environment...")

	return nil
}
