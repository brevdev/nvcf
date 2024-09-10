package deployment

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/brevdev/nvcf/api"
	"github.com/brevdev/nvcf/config"
	"github.com/brevdev/nvcf/output"
	"github.com/spf13/cobra"
)

func deploymentGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <function-id>",
		Short: "Get a NVIDIA Cloud Deployment. If a version-id is not provided and there are multiple associated with a function, we will look for all deployments and prompt for a version-id.",
		Args:  cobra.ExactArgs(1),
		Run:   runDeploymentGet,
	}

	cmd.Flags().String("version-id", "", "The version ID of the deployment")
	return cmd
}

func runDeploymentGet(cmd *cobra.Command, args []string) {
	client := api.NewClient(config.GetAPIKey())
	versionID, _ := cmd.Flags().GetString("version-id")
	// version id was provided
	if versionID != "" {
		deployment, err := client.FunctionDeployment.Functions.Versions.GetDeployment(cmd.Context(), args[0], versionID)
		if err != nil {
			output.Error(cmd, "Error getting deployment", err)
			return
		}
		fmt.Println(deployment)
		// output.Deployment(cmd, deployment)
	}

	// version id was not provided - lets look for the function deployment
	versions, err := client.Functions.Versions.List(cmd.Context(), args[0])
	if err != nil {
		output.Error(cmd, "Error listing function versions", err)
		return
	}
	// there is only 1 version which is the one that was deployed
	if len(versions.Functions) == 1 {
		deployment, err := client.FunctionDeployment.Functions.Versions.GetDeployment(cmd.Context(), args[0], versions.Functions[0].ID)
		if err != nil {
			output.Error(cmd, "Error getting deployment", err)
			return
		}
		fmt.Println(deployment)
		// output.SingleFunction(cmd, deployment)
	} else {
		output.Info(cmd, "Multiple versions found. Please specify a version-id")
		for _, version := range versions.Functions {
			// skip functions that have a status of INACTIVE
			if version.Status == "INACTIVE" {
				continue
			}
			output.Info(cmd, fmt.Sprintf("Version ID: %s || Status: %s", version.VersionID, version.Status))
		}
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter version-id: ")
		versionID, _ = reader.ReadString('\n')
		versionID = strings.TrimSpace(versionID)
		deployment, err := client.FunctionDeployment.Functions.Versions.GetDeployment(cmd.Context(), args[0], versionID)
		if err != nil {
			output.Error(cmd, "Error getting deployment", err)
			return
		}
		fmt.Println(deployment)
		// output.Deployment(cmd, deployment)
	}
}