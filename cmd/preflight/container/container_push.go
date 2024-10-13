package container

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/brevdev/nvcf/cmd/auth"
	"github.com/spf13/cobra"
)

func PushCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push",
		Short: "Tag a new container and push it to nvcr.io",
		Long:  "Tag a new container and push it to nvcr.io",
		Args:  cobra.ExactArgs(1),
		RunE:  runPush,
	}

	cmd.Flags().StringP("tag", "t", "latest", "Tag to push")

	return cmd
}

func runPush(cmd *cobra.Command, args []string) error {
	image := args[0]
	tag, err := cmd.Flags().GetString("tag")
	if err != nil {
		return err
	}

	// Check if the image exists locally
	if err := checkAndPullImage(image, tag); err != nil {
		return err
	}

	// Get the orgId
	orgId, err := auth.GetOrgId()
	if err != nil {
		return fmt.Errorf("failed to get organization ID: %w", err)
	}

	// Split the image name to separate the repository and image name
	imageParts := strings.Split(image, "/")
	imageName := imageParts[len(imageParts)-1]

	// Construct the new image name
	newImageName := fmt.Sprintf("nvcr.io/%s/%s", orgId, imageName)

	// Retag the image
	retagCmd := exec.Command("docker", "tag", fmt.Sprintf("%s:%s", image, tag), fmt.Sprintf("%s:%s", newImageName, tag))
	if output, err := retagCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to retag image: %s, %w", string(output), err)
	}

	fmt.Printf("Successfully retagged image as %s:%s\n", newImageName, tag)

	// Push the new image
	pushCmd := exec.Command("docker", "push", fmt.Sprintf("%s:%s", newImageName, tag))

	// Set up pipes to capture and display output in real-time
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr

	if err := pushCmd.Run(); err != nil {
		return fmt.Errorf("failed to push image: %w", err)
	}

	fmt.Printf("Successfully pushed image %s:%s to nvcr.io\n", newImageName, tag)

	return nil
}

func checkAndPullImage(image, tag string) error {
	// Check if the image exists locally
	inspectCmd := exec.Command("docker", "inspect", fmt.Sprintf("%s:%s", image, tag))
	if err := inspectCmd.Run(); err == nil {
		// Image exists locally
		fmt.Printf("Image %s:%s found locally\n", image, tag)
		return nil
	}

	// Image doesn't exist locally, pull it
	fmt.Printf("Image %s:%s not found locally. Pulling...\n", image, tag)
	pullCmd := exec.Command("docker", "pull", fmt.Sprintf("%s:%s", image, tag))
	pullCmd.Stdout = os.Stdout
	pullCmd.Stderr = os.Stderr

	if err := pullCmd.Run(); err != nil {
		return fmt.Errorf("failed to pull image %s:%s: %w", image, tag, err)
	}

	fmt.Printf("Successfully pulled image %s:%s\n", image, tag)
	return nil
}
