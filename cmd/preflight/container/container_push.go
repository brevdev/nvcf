package container

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/brevdev/nvcf/api"
	"github.com/brevdev/nvcf/config"
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

	// Get the orgId
	orgId, err := getOrgId()
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

// getOrgId retrieves the organization ID
func getOrgId() (string, error) {
	client := api.NewClient(config.GetAPIKey())
	orgsInfo := map[string]interface{}{}
	err := client.Get(context.Background(), "/v2/orgs", nil, &orgsInfo)
	if err != nil {
		return "", err
	}

	organizations, ok := orgsInfo["organizations"].([]interface{})
	if !ok || len(organizations) == 0 {
		return "", fmt.Errorf("no organizations found")
	}

	firstOrg, ok := organizations[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("failed to parse organization information")
	}

	name, ok := firstOrg["name"].(string)
	if !ok {
		return "", fmt.Errorf("organization name not found")
	}

	return name, nil
}
