package containerutil

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"time"
)

type ContainerSmokeTest struct {
	DefaultPort    int
	WaitIterations int
	ContainerID    string
	HostPort       string
}

func NewContainerSmokeTest() (*ContainerSmokeTest, error) {
	return &ContainerSmokeTest{
		DefaultPort:    18080,
		WaitIterations: 100,
	}, nil
}

func (cst *ContainerSmokeTest) LaunchContainer(imageName, containerPort string) error {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "docker", "run", "-d", "-p", fmt.Sprintf("%d:%s", cst.DefaultPort, containerPort), imageName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start container: %v, output: %s", err, output)
	}

	cst.ContainerID = string(bytes.TrimSpace(output))
	cst.HostPort = fmt.Sprintf("%d", cst.DefaultPort)

	// Check if the container is still running after a short delay
	time.Sleep(2 * time.Second)

	cmd = exec.CommandContext(ctx, "docker", "inspect", "-f", "{{.State.Running}}", cst.ContainerID)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to inspect container: %v, output: %s", err, output)
	}

	if string(bytes.TrimSpace(output)) != "true" {
		// Container exited, fetch and display logs
		cmd = exec.CommandContext(ctx, "docker", "logs", cst.ContainerID)
		logs, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("container exited and failed to fetch logs: %v", err)
		}

		return fmt.Errorf("container exited unexpectedly. Logs:\n%s", string(logs))
	}

	return nil
}

func (cst *ContainerSmokeTest) CheckHTTPHealthEndpoint(healthEndpoint string, secondsToWaitForHealthy int) error {
	healthURL := fmt.Sprintf("http://localhost:%s%s", cst.HostPort, healthEndpoint)
	fmt.Printf("Looking for health signal at %s\n", healthURL)

	client := &http.Client{Timeout: 5 * time.Second}

	for i := 0; i < cst.WaitIterations; i++ {
		resp, err := client.Get(healthURL)
		if err == nil && resp.StatusCode == 200 {
			fmt.Printf("%s returned 200 OK\n", healthURL)
			return nil
		}
		time.Sleep(time.Duration(secondsToWaitForHealthy/cst.WaitIterations) * time.Second)
	}

	return fmt.Errorf("health check did not complete successfully in time")
}

func (cst *ContainerSmokeTest) Cleanup() error {
	cmd := exec.Command("docker", "rm", "-f", cst.ContainerID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove container: %v, output: %s", err, output)
	}
	return nil
}

func (cst *ContainerSmokeTest) ForceCleanup(imageName string) error {
	cmd := exec.Command("docker", "ps", "-a", "-q", "--filter", fmt.Sprintf("ancestor=%s", imageName))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to list containers: %v, output: %s", err, output)
	}

	containerIDs := bytes.Fields(output)
	for _, containerID := range containerIDs {
		stopCmd := exec.Command("docker", "stop", string(containerID))
		_, err := stopCmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Warning: Failed to stop container %s: %v\n", containerID, err)
		}

		rmCmd := exec.Command("docker", "rm", "-f", string(containerID))
		_, err = rmCmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Warning: Failed to remove container %s: %v\n", containerID, err)
		}
	}

	return nil
}

func (cst *ContainerSmokeTest) TestHTTPInference(inferenceEndpoint string, payload interface{}) error {
	inferenceURL := fmt.Sprintf("http://localhost:%s%s", cst.HostPort, inferenceEndpoint)
	fmt.Printf("Sending payload to %s...\n", inferenceURL)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	resp, err := http.Post(inferenceURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to send inference request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("Server's response status: %s\n", resp.Status)
	fmt.Printf("Server's response headers: %v\n", resp.Header)

	if resp.Header.Get("Content-Type") == "text/event-stream" {
		fmt.Println("Received a streaming response:")
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading stream: %v", err)
		}
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %v", err)
		}
		fmt.Printf("Server's response body: %s\n", string(body))

		var result interface{}
		if err := json.Unmarshal(body, &result); err == nil {
			prettyJSON, _ := json.MarshalIndent(result, "", "  ")
			fmt.Printf("Parsed JSON response:\n%s\n", string(prettyJSON))
		}
	}

	return nil
}
