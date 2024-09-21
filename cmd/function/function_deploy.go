package function

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/brevdev/nvcf/api"
	"github.com/brevdev/nvcf/config"
	"github.com/brevdev/nvcf/output"
	"github.com/spf13/cobra"
	"github.com/tmc/nvcf-go"
)

func functionDeployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deploy <function-id>",
		Short:   "Deploy a function",
		Long:    `Deploy an existing NVCF function. If you want to deploy a specific version, use the --version-id flag.`,
		Example: "nvcf function deploy fid --version-id vid --gpu A100 --instance-type g5.4xlarge",
		Args:    cobra.ExactArgs(1),
		RunE:    runFunctionDeploy,
	}

	cmd.Flags().String("version-id", "", "The ID of the version to deploy")
	cmd.Flags().String("gpu", "", "GPU type to use")
	cmd.Flags().String("instance-type", "", "Instance type to use")
	cmd.Flags().String("backend", "", "Backend to deploy the function to")
	cmd.Flags().Int64("min-instances", 0, "Minimum number of instances")
	cmd.Flags().Int64("max-instances", 1, "Maximum number of instances")
	cmd.Flags().Int64("max-request-concurrency", 1, "Maximum number of concurrent requests")
	cmd.Flags().BoolP("detached", "d", false, "Detach from the deployment and return to the prompt")

	cmd.MarkFlagRequired("gpu")
	cmd.MarkFlagRequired("instance-type")
	cmd.MarkFlagRequired("backend")

	return cmd
}

func runFunctionDeploy(cmd *cobra.Command, args []string) error {
	client := api.NewClient(config.GetAPIKey())
	functionId := args[0]
	versionId, _ := cmd.Flags().GetString("version-id")

	if versionId == "" {
		versions, err := client.Functions.Versions.List(cmd.Context(), functionId)
		if err != nil {
			return output.Error(cmd, "Error listing function versions", err)
		}

		if len(versions.Functions) == 1 {
			versionId = versions.Functions[0].VersionID
		} else {
			output.Info(cmd, "Multiple versions found. Please specify a version-id")
			for _, version := range versions.Functions {
				output.Info(cmd, fmt.Sprintf("Version ID: %s || Status: %s", version.VersionID, version.Status))
			}
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter version-id: ")
			versionId, _ = reader.ReadString('\n')
			versionId = strings.TrimSpace(versionId)
		}
	}

	gpu, _ := cmd.Flags().GetString("gpu")
	instanceType, _ := cmd.Flags().GetString("instance-type")
	backend, _ := cmd.Flags().GetString("backend")
	minInstances, _ := cmd.Flags().GetInt64("min-instances")
	maxInstances, _ := cmd.Flags().GetInt64("max-instances")
	maxRequestConcurrency, _ := cmd.Flags().GetInt64("max-request-concurrency")
	detached, _ := cmd.Flags().GetBool("detached")

	deploymentParams := nvcf.FunctionDeploymentFunctionVersionInitiateDeploymentParams{
		DeploymentSpecifications: nvcf.F([]nvcf.FunctionDeploymentFunctionVersionInitiateDeploymentParamsDeploymentSpecification{{
			GPU:                   nvcf.String(gpu),
			InstanceType:          nvcf.String(instanceType),
			Backend:               nvcf.String(backend),
			MaxInstances:          nvcf.Int(maxInstances),
			MinInstances:          nvcf.Int(minInstances),
			MaxRequestConcurrency: nvcf.Int(maxRequestConcurrency),
		}}),
	}

	_, err := client.FunctionDeployment.Functions.Versions.InitiateDeployment(
		cmd.Context(),
		functionId,
		versionId,
		deploymentParams,
	)
	if err != nil {
		return output.Error(cmd, "Error deploying function", err)
	}

	output.Success(cmd, fmt.Sprintf("Function %s version %s deployed successfully", functionId, versionId))
	if !detached {
		return WaitForDeployment(cmd, client, functionId, versionId)
	}

	return nil
}
