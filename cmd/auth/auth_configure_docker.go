package auth

import (
	"os/exec"
	"strings"

	"github.com/brevdev/nvcf/config"
	"github.com/brevdev/nvcf/output"
	"github.com/spf13/cobra"
)

func authConfigureDockerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "configure-docker",
		Short: "Configure Docker to use NGC API key for nvcr.io",
		RunE: func(cmd *cobra.Command, args []string) error {
			apiKey := config.GetAPIKey()
			if apiKey == "" {
				return output.Error(cmd, "NGC API key not found. Please run 'nvcf auth login' first.", nil)
			}
			// Check if Docker is installed
			_, err := exec.LookPath("docker")
			if err != nil {
				return output.Error(cmd, "Docker is not installed or not in the system PATH", err)
			}
			// TODO: check for existing nvcr.io config?
			dockerCmd := exec.Command("docker", "login", "nvcr.io", "-u", "$oauthtoken", "--password-stdin")
			dockerCmd.Stdin = strings.NewReader(apiKey)
			out, err := dockerCmd.CombinedOutput()
			if err != nil {
				cmd.Println(string(out))
				return output.Error(cmd, "Failed to configure Docker", err)
			}
			output.Success(cmd, "Docker configured successfully for nvcr.io")
			cmd.Println(string(out))
			return nil
		},
	}
}
