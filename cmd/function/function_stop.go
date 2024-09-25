package function

import (
	"fmt"

	"github.com/brevdev/nvcf/api"
	"github.com/brevdev/nvcf/config"
	"github.com/brevdev/nvcf/output"
	"github.com/spf13/cobra"
	"github.com/tmc/nvcf-go"
)

func functionStopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stop <function-id>",
		Short:   "Stop a deployed function",
		Long:    "Stop a deployed function. If a version-id is not provided, we look for versions that are actively deployed. If a single function is deployed, we stop that version. If multiple functions are deployed, we prompt for the version-id to stop.",
		Example: "nvcf function stop fid --version-id vid",
		Args:    cobra.ExactArgs(1),
		RunE:    runFunctionStop,
	}
	cmd.Flags().Bool("all", false, "Stop all deployed versions of the function")
	cmd.Flags().String("version-id", "", "The ID of the version")
	cmd.Flags().Bool("force", false, "Gracefully stop the function if it's already deployed. If not, we forcefully stop the function.")
	return cmd
}

func runFunctionStop(cmd *cobra.Command, args []string) error {
	client := api.NewClient(config.GetAPIKey())

	functionId := args[0]
	versionId, _ := cmd.Flags().GetString("version-id")
	force, _ := cmd.Flags().GetBool("force")
	all, _ := cmd.Flags().GetBool("all")

	if versionId == "" {
		versions, err := client.Functions.Versions.List(cmd.Context(), functionId)
		if err != nil {
			return output.Error(cmd, "Error listing function versions", err)
		}
		deployedVersionsToStop := []string{}
		for _, version := range versions.Functions {
			if version.Status == "ACTIVE" {
				deployedVersionsToStop = append(deployedVersionsToStop, version.VersionID)
			}
		}
		if len(deployedVersionsToStop) == 0 {
			return output.Error(cmd, "No functions are currently deployed", nil)
		}
		if all {
			for _, version := range deployedVersionsToStop {
				client.FunctionDeployment.Functions.Versions.DeleteDeployment(cmd.Context(), functionId, version, nvcf.FunctionDeploymentFunctionVersionDeleteDeploymentParams{
					Graceful: nvcf.Bool(force),
				})
				output.Success(cmd, fmt.Sprintf("Function %s version %s stopped successfully", functionId, version))
			}
		} else {
			client.FunctionDeployment.Functions.Versions.DeleteDeployment(cmd.Context(), functionId, deployedVersionsToStop[0], nvcf.FunctionDeploymentFunctionVersionDeleteDeploymentParams{
				Graceful: nvcf.Bool(force),
			})
			output.Success(cmd, fmt.Sprintf("Function %s version %s stopped successfully", functionId, deployedVersionsToStop[0]))
		}
	}
	return nil
}
