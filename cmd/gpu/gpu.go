package gpu

import (
	"fmt"
	"os"

	"github.com/brevdev/nvcf/config"
	"github.com/spf13/cobra"
)

func GpuCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "gpu",
		Aliases: []string{"gpus", "acc", "accelerators", "accelerator-types"},
		Short:   "Manage cluster groups and available GPUs",
		Long:    `List available GPUs, cluster groups, and other GPU related information`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			config.Init()
			if cmd.Name() != "auth" && !config.IsAuthenticated() {
				fmt.Println("You are not authenticated. Please run 'nvcf auth login' first.")
				os.Exit(1)
			}
		},
	}

	cmd.AddCommand(gpuListCmd())

	return cmd
}
