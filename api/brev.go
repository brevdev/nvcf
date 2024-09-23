package api

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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

func (c *BrevClient) CreateInstance(instanceName string) (string, error) {
	cmd := exec.Command("brev", "create", instanceName, "-g", instanceType)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to create instance: %w", err)
	}
	return instanceName, nil
}
