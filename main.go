package main

import (
	"fmt"
	"os"

	"github.com/brevdev/nvcf/cmd"
	"github.com/brevdev/nvcf/cmd/auth"
	"github.com/brevdev/nvcf/cmd/function"
	"github.com/brevdev/nvcf/cmd/gpu"
	"github.com/brevdev/nvcf/cmd/preflight"
	"github.com/brevdev/nvcf/config"
	"github.com/brevdev/nvcf/output"
	"github.com/spf13/cobra"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	rootCmd := &cobra.Command{
		Use:   "nvcf",
		Short: "NVIDIA Cloud Functions CLI",
		Long: `A command-line interface for managing and interacting with NVIDIA Cloud Functions.

Environment variables:
  NVCF_BETA - Set to true to enable beta features
  NVCF_SHOW_DOCS_CMD - Set to true to show the docs command
`,
		SilenceErrors:     true,
		PersistentPreRunE: preRunAuthCheck,
		Run: func(cmd *cobra.Command, args []string) {
			output.PrintASCIIArt(cmd)
			err := cmd.Usage()
			if err != nil {
				return
			}
		},
		DisableAutoGenTag: true,
	}

	// Add global flags
	rootCmd.PersistentFlags().Bool("json", false, "Output results in JSON format")
	rootCmd.PersistentFlags().Bool("no-color", false, "Disable color output")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Suppress non-error output")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output and show underlying API calls")

	// Add commands
	rootCmd.AddCommand(function.FunctionCmd())
	rootCmd.AddCommand(gpu.GpuCmd())
	// rootCmd.AddCommand(cmd.InvokeCmd())
	// rootCmd.AddCommand(cmd.AssetCmd())
	rootCmd.AddCommand(auth.AuthCmd())
	// rootCmd.AddCommand(cmd.QueueCmd())
	// rootCmd.AddCommand(cmd.ClusterGroupCmd())
	// rootCmd.AddCommand(cmd.ConfigCmd())
	rootCmd.AddCommand(preflight.PreflightCmd())
	rootCmd.AddCommand(cmd.DocsCmd())

	// // Enable command auto-completion
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(cmd.CompletionCmd())

	return rootCmd.Execute()
}

var shouldApplyAuthCheck = map[string]bool{
	"function": true,
	"gpu":      true,
}

func preRunAuthCheck(cmd *cobra.Command, args []string) error {
	config.Init()
	topLevelCmd := getTopLevelCmd(cmd)
	if shouldApplyAuthCheck[topLevelCmd.Name()] {
		if !config.IsAuthenticated() {
			return fmt.Errorf("you are not authenticated. Please run 'nvcf auth login' first")
		}
	}
	return nil
}

func getTopLevelCmd(cmd *cobra.Command) *cobra.Command {
	if cmd == nil || cmd.Parent() == nil {
		return cmd
	}
	parent := cmd.Parent()
	if parent.Parent() == nil {
		return cmd
	}
	return getTopLevelCmd(parent)
}
