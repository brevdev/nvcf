package debug

import (
	"os"

	"github.com/spf13/cobra"
)

func DebugCmd() *cobra.Command {
	cmd := &cobra.Command{
		Aliases: []string{"d", "dbg"},
		Use:     "debug",
		Short:   "Debug NVCF functions",
		Long:    `Create and manage debug environments for NVCF functions`,
		Hidden:  os.Getenv("NVCF_BETA") != "true",
	}

	cmd.AddCommand(debugStartCmd())
	cmd.AddCommand(debugStopCmd())

	return cmd
}
