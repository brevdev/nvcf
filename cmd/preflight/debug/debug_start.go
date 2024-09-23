package debug

import (
	"fmt"

	"github.com/brevdev/nvcf/api"
	"github.com/spf13/cobra"
)

func debugStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start a debug environment",
		Long:  `Create and start a debug environment for an NVCF function`,
		RunE:  runDebugStart,
	}

	cmd.Flags().StringP("instance-name", "n", "", "The name of the instance to debug")
	return cmd
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
	fmt.Println("Setting up a GPU powered VM for debugging")

	instanceName, _ := cmd.Flags().GetString("instance-name")

	if instanceName == "" {
		return fmt.Errorf("instance name is required. Please provide an instance name")
	}

	// hit the brev api to create an instance using
	brevClient.CreateInstance(instanceName)

	fmt.Sprintf("you can enter this instance for debugging purposes using ssh %s-host", instanceName)

	return nil
}
