package preflight

import (
	"github.com/brevdev/nvcf/cmd/preflight/check"
	"github.com/brevdev/nvcf/cmd/preflight/debug"
	"github.com/brevdev/nvcf/config"
	"github.com/spf13/cobra"
)

func PreflightCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "preflight",
		Short: "Perform preflight checks for NVCF compatibility",
		Long:  "Run various preflight checks to ensure compatibility with NVIDIA Cloud Functions (NVCF).",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			config.Init()
			return nil
		},
	}

	cmd.AddCommand(check.CheckCmd())
	cmd.AddCommand(debug.DebugCmd())
	return cmd
}
