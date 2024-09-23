package preflight

import (
	"github.com/brevdev/nvcf/cmd/preflight/check"
	"github.com/spf13/cobra"
)

func PreflightCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "preflight",
		Short: "Perform preflight checks for NVCF compatibility",
		Long:  "Run various preflight checks to ensure compatibility with NVIDIA Cloud Functions (NVCF).",
	}

	cmd.AddCommand(check.NewCheckCmd())

	return cmd
}
