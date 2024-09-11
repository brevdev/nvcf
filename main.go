package main

import (
	"fmt"
	"os"

	"github.com/brevdev/nvcf/cmd"
	"github.com/brevdev/nvcf/cmd/auth"
	"github.com/brevdev/nvcf/cmd/function"
	"github.com/brevdev/nvcf/cmd/gpu"
	"github.com/brevdev/nvcf/cmd/test"
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
		Use:           "nvcf",
		Short:         "NVIDIA Cloud Functions CLI",
		Long:          `A command-line interface for managing and interacting with NVIDIA Cloud Functions.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(cmd *cobra.Command, args []string) {
			output.PrintASCIIArt(cmd)
			cmd.Usage()
		},
	}

	// Add global flags
	rootCmd.PersistentFlags().Bool("json", false, "Output results in JSON format")
	rootCmd.PersistentFlags().Bool("no-color", false, "Disable color output")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Suppress non-error output")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")

	// Add commands
	rootCmd.AddCommand(function.FunctionCmd())
	rootCmd.AddCommand(gpu.GpuCmd())
	// rootCmd.AddCommand(cmd.InvokeCmd())
	// rootCmd.AddCommand(cmd.AssetCmd())
	rootCmd.AddCommand(auth.AuthCmd())
	// rootCmd.AddCommand(cmd.QueueCmd())
	// rootCmd.AddCommand(cmd.ClusterGroupCmd())
	// rootCmd.AddCommand(cmd.ConfigCmd())
	rootCmd.AddCommand(cmd.DocsCmd())
	rootCmd.AddCommand(test.TestCmd())

	// // Enable command auto-completion
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(cmd.CompletionCmd())

	return rootCmd.Execute()
}
