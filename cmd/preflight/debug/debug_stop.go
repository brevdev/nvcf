package debug

import (
	"fmt"

	"github.com/brevdev/nvcf/api"
	"github.com/spf13/cobra"
)

func debugStopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop a debug environment",
		Long:  `Stop and clean up a debug environment for an NVCF function`,
		RunE:  runDebugStop,
	}
}

func runDebugStop(cmd *cobra.Command, args []string) error {
	brevClient := api.NewBrevClient()

	if !brevClient.IsBrevCLIInstalled() {
		return fmt.Errorf("brev CLI is not installed. Please install it first")
	}

	instanceName, _ := cmd.Flags().GetString("instance-name")

	if instanceName == "" {
		return fmt.Errorf("instance name is required. Please provide an instance name")
	}

	err := brevClient.DeleteInstance(instanceName)
	if err != nil {
		return fmt.Errorf("error deleting instance: %w", err)
	}

	return nil
}
