package cmd

import (
	"os/exec"
	"strings"

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
	// cmd.AddCommand(authLogoutCmd())
	// cmd.AddCommand(authStatusCmd())

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

func authConfigureDockerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "configure-docker",
		Short: "Configure Docker to use NGC API key for nvcr.io",
		Run: func(cmd *cobra.Command, args []string) {
			apiKey := config.GetAPIKey()
			if apiKey == "" {
				output.Error(cmd, "NGC API key not found. Please run 'nvcf auth login' first.", nil)
				return
			}
			// TODO: check for 'docker'
			// TODO: check for existing nvcr.io config?
			dockerCmd := exec.Command("docker", "login", "nvcr.io", "-u", "$oauthtoken", "--password-stdin")
			dockerCmd.Stdin = strings.NewReader(apiKey)
			out, err := dockerCmd.CombinedOutput()
			if err != nil {
				output.Error(cmd, "Failed to configure Docker", err)
				cmd.Println(string(out))
				return
			}
			output.Success(cmd, "Docker configured successfully for nvcr.io")
			cmd.Println(string(out))
		},
	}
}

// Implement authLogoutCmd and authStatusCmd here
