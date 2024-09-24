package api

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/brevdev/nvcf/config"
)

type BrevClient struct{}

func NewBrevClient() *BrevClient {
	return &BrevClient{}
}

func (c *BrevClient) IsBrevCLIInstalled() bool {
	_, err := exec.LookPath("brev")
	return err == nil
}

func (c *BrevClient) Login() error {
	cmd := exec.Command("brev", "login")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to login with Brev CLI: %w", err)
	}
	return nil
}

func (c *BrevClient) IsLoggedIn() (bool, error) {
	cmd := exec.Command("brev", "ls")
	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "currently logged out") {
			return false, nil
		}
		return false, fmt.Errorf("failed to check login status: %w\nOutput: %s", err, output)
	}
	return true, nil
}

const (
	instanceType = "n1-standard-8:nvidia-tesla-t4:1"
)

func (c *BrevClient) CreateInstance(instanceName string) error {
	cmd := exec.Command("brev", "create", instanceName, "-g", instanceType)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to create instance: %w", err)
	}
	return nil
}

func (c *BrevClient) RunDebuggingScript(instanceName string, image string, imageArgs string) error {
	debuggingScript := generateDebuggingScript(image, imageArgs)

	// Escape single quotes in the script
	escapedScript := strings.ReplaceAll(debuggingScript, "'", "'\\''")

	cmd := exec.Command("brev", "refresh")
	cmd.Run()

	sshAlias := instanceName
	sshCmd := []string{
		"bash",
		"-c",
		escapedScript,
	}

	err := runSSHExec(sshAlias, sshCmd, false)
	if err != nil {
		return fmt.Errorf("failed to run debugging script: %w", err)
	}
	return nil
}

func (c *BrevClient) DeleteInstance(instanceName string) error {
	cmd := exec.Command("brev", "delete", instanceName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to delete instance: %w", err)
	}
	return nil
}

func runSSHExec(sshAlias string, args []string, fireAndForget bool) error {
	sshAgentEval := "eval $(ssh-agent -s)"
	cmd := fmt.Sprintf("ssh %s -- %s", sshAlias, strings.Join(args, " "))
	fmt.Printf("cmd1: %s\n", cmd)
	cmd = fmt.Sprintf("%s && %s", sshAgentEval, cmd)
	fmt.Printf("cmd2: %s\n", cmd)
	sshCmd := exec.Command("bash", "-c", cmd) //nolint:gosec //cmd is user input

	if fireAndForget {
		if err := sshCmd.Start(); err != nil {
			return fmt.Errorf("error starting SSH command: %w", err)
		}
		return nil
	}
	sshCmd.Stderr = os.Stderr
	sshCmd.Stdout = os.Stdout
	sshCmd.Stdin = os.Stdin
	if err := sshCmd.Run(); err != nil {
		return fmt.Errorf("error running SSH command: %w", err)
	}
	return nil
}

func generateDebuggingScript(image string, imageArgs string) string {
	return fmt.Sprintf(`
#!/bin/bash

# Start the debugging session
echo "Starting debugging session"

# Install dependencies
echo "Logging into nvcr.io using API credentials"
echo %s | docker login nvcr.io --username $oauthtoken --password-stdin

# Pull the container image
docker pull %s

# Run the container image
docker run -it %s

echo "Debugging session complete"
`, config.GetAPIKey(), image, imageArgs)
}
