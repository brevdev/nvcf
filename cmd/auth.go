package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tmc/nvcf/config"
	"github.com/tmc/nvcf/output"
)

func AuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication for the CLI",
	}

	cmd.AddCommand(authLoginCmd())
	cmd.AddCommand(authLogoutCmd())
	cmd.AddCommand(authStatusCmd())

	return cmd
}

func authLoginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Authenticate with NVIDIA Cloud",
		Run: func(cmd *cobra.Command, args []string) {
			apiKey := output.Prompt("Enter your NVIDIA Cloud API key: ", true)
			err := config.SetAPIKey(apiKey)
			if err != nil {
				output.Error(cmd, "Error saving API key", err)
				return
			}
			output.Success(cmd, "Authentication successful")
		},
	}
}

// Implement authLogoutCmd and authStatusCmd here
