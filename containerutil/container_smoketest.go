package containerutil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type ContainerSmokeTest struct {
	DefaultPort    int
	WaitIterations int
	DockerClient   *client.Client
	ContainerID    string
	HostPort       string
}

func NewContainerSmokeTest() (*ContainerSmokeTest, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %v", err)
	}

	return &ContainerSmokeTest{
		DefaultPort:    18080,
		WaitIterations: 100,
		DockerClient:   cli,
	}, nil
}

func (cst *ContainerSmokeTest) LaunchContainer(imageName, containerPort string) error {
	ctx := context.Background()
	resp, err := cst.DockerClient.ContainerCreate(ctx,
		&container.Config{
			Image: imageName,
			ExposedPorts: nat.PortSet{
				nat.Port(containerPort): struct{}{},
			},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				nat.Port(containerPort): []nat.PortBinding{
					{
						HostIP:   "127.0.0.1",
						HostPort: fmt.Sprintf("%d", cst.DefaultPort),
					},
				},
			},
		}, nil, nil, "")
	if err != nil {
		return fmt.Errorf("failed to create container: %v", err)
	}

	cst.ContainerID = resp.ID
	cst.HostPort = fmt.Sprintf("%d", cst.DefaultPort)

	if err := cst.DockerClient.ContainerStart(ctx, cst.ContainerID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %v", err)
	}

	// Check if the container is still running after a short delay
	time.Sleep(2 * time.Second)

	containerInfo, err := cst.DockerClient.ContainerInspect(ctx, cst.ContainerID)
	if err != nil {
		return fmt.Errorf("failed to inspect container: %v", err)
	}

	if !containerInfo.State.Running {
		// Container exited, fetch and display logs
		out, err := cst.DockerClient.ContainerLogs(ctx, cst.ContainerID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
		if err != nil {
			return fmt.Errorf("container exited and failed to fetch logs: %v", err)
		}
		defer out.Close()

		logs, err := io.ReadAll(out)
		if err != nil {
			return fmt.Errorf("failed to read container logs: %v", err)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	fmt.Printf("Server's response: %s\n", string(body))
	return nil
}

func (cst *ContainerSmokeTest) Cleanup() error {
	ctx := context.Background()
	return cst.DockerClient.ContainerRemove(ctx, cst.ContainerID, container.RemoveOptions{Force: true})
}

// ForceCleanup stops and removes all containers created from the specified image
func (cst *ContainerSmokeTest) ForceCleanup(imageName string) error {
	ctx := context.Background()
	containers, err := cst.DockerClient.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return fmt.Errorf("failed to list containers: %v", err)
	}

	stopOpts := container.StopOptions{}
	removeOpts := container.RemoveOptions{Force: true}
	for _, container := range containers {
		if container.Image == imageName {
			fmt.Printf("Stopping and removing container %s...\n", container.ID[:12])
			if err := cst.DockerClient.ContainerStop(ctx, container.ID, stopOpts); err != nil {
				fmt.Printf("Warning: Failed to stop container %s: %v\n", container.ID[:12], err)
			}
			if err := cst.DockerClient.ContainerRemove(ctx, container.ID, removeOpts); err != nil {
				fmt.Printf("Warning: Failed to remove container %s: %v\n", container.ID[:12], err)
			}
		}
	}
	return nil
}
