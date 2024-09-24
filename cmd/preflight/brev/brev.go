package brev

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

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
	instanceType = "n1-standard-8:nvidia-tesla-t4:1" // gcp t4
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

	cmd := exec.Command("brev", "refresh")
	cmd.Run()

	sshAlias := instanceName
	sshCmd := []string{
		debuggingScript,
	}

	// Retry SSH connection
	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		err = runSSHExec(sshAlias, sshCmd)
		if err == nil {
			return nil
		}
		fmt.Printf("Connecting to instance %s\n", instanceName)
		// fmt.Printf("SSH connection attempt %d failed: %v. Retrying...\n", i+1, err)
		time.Sleep(7 * time.Second)
	}
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

func runSSHExec(sshAlias string, args []string) error {
	sshCmd := exec.Command("ssh", sshAlias)
	si, err := sshCmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("error getting stdin pipe: %w", err)
	}
	for _, arg := range args {
		si.Write([]byte(arg + "\n"))
	}
	si.Close()

	sshCmd.Stderr = os.Stderr
	sshCmd.Stdout = os.Stdout
	if err := sshCmd.Run(); err != nil {
		return fmt.Errorf("Connecting")
	}
	return nil
}

func generateDebuggingScript(image string, imageArgs string) string {
	return fmt.Sprintf(`
# Start the debugging session
echo "Starting debugging session"

sudo docker ps || true

# Install dependencies
echo "Logging into nvcr.io using API credentials"
echo %s | sudo docker login nvcr.io --username '$oauthtoken' --password-stdin || true

# Pull the container image
echo "Pulling image %s"
sudo docker pull %s --quiet

# Run the container image
sudo docker run %s %s

echo "Debugging session complete"
`, config.GetAPIKey(), image, image, image, imageArgs)
}
