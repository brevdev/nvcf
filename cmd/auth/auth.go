package auth

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/brevdev/nvcf/config"
	"github.com/brevdev/nvcf/output"
	"github.com/spf13/cobra"
)

func AuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication for the CLI",
		Long:  `Authenticate with NVIDIA Cloud and configure the CLI to use your API key.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			config.Init()
			// Allow all 'auth' subcommands to run without authentication
			if cmd.Parent().Name() != "auth" && !config.IsAuthenticated() {
				fmt.Println("You are not authenticated. Please run 'nvcf auth login' first.")
				os.Exit(1)
			}
		},
	}

	cmd.AddCommand(authLoginCmd())
	cmd.AddCommand(authLogoutCmd())
	cmd.AddCommand(authStatusCmd())
	cmd.AddCommand(authConfigureDockerCmd())

	return cmd
}

func authLoginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Authenticate with NVIDIA Cloud",
		Run: func(cmd *cobra.Command, args []string) {
			apiKey := output.Prompt("Enter your NVIDIA Cloud API key: ", true)
			orgID := output.Prompt("Enter your NVIDIA Cloud Org ID: ", true)
			err := config.SetAPIKey(apiKey)
			if err != nil {
				output.Error(cmd, "Error saving API key", err)
				return
			}
			err = config.SetOrgID(orgID)
			if err != nil {
				output.Error(cmd, "Error saving Org ID", err)
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
			// Check if Docker is installed
			_, err := exec.LookPath("docker")
			if err != nil {
				output.Error(cmd, "Docker is not installed or not in the system PATH", err)
				return
			}
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

func authStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Check the authentication status",
		Run: func(cmd *cobra.Command, args []string) {
			if config.IsAuthenticated() {
				output.Success(cmd, "Authenticated")
			} else {
				output.Error(cmd, "Not authenticated", errors.New(":("))
			}
		},
	}
}

func authLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Logout from NVIDIA Cloud",
		Run: func(cmd *cobra.Command, args []string) {
			if !config.IsAuthenticated() {
				output.Info(cmd, "You are currently not logged in")
				return
			}
			config.ClearAPIKey()
			output.Success(cmd, "Logged out successfully")
		},
	}
}
