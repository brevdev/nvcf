package debug

import (
	"github.com/spf13/cobra"
)

func DebugCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "debug",
		Short: "Debug NVCF functions",
		Long:  `Create and manage debug environments for NVCF functions`,
	}

	cmd.AddCommand(debugStartCmd())
	cmd.AddCommand(debugStopCmd())

	return cmd
}
