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
		// check if this error is due to an ongoing deployment
		// get the function and check the status
		fn, err := client.Functions.Versions.Get(cmd.Context(), functionId, versionId, nvcf.FunctionVersionGetParams{
			IncludeSecrets: nvcf.Bool(false),
		})
		if err != nil {
			return output.Error(cmd, "Error checking function status. Please try again", err)
		}
		if fn.Function.Status != nvcf.FunctionResponseFunctionStatusInactive {
			output.Info(cmd, fmt.Sprintf("This function is currently %s. ", fn.Function.Status))
			output.Info(cmd, "Creating new version and deploying...")
			newVersionToDeploy, err := createNewVersion(cmd, client, fn.Function)
			if err != nil {
				return output.Error(cmd, "Error creating new version to deploy", err)
			}
			_, err = client.FunctionDeployment.Functions.Versions.InitiateDeployment(
				cmd.Context(),
				functionId,
				newVersionToDeploy,
				deploymentParams,
			)
			if err != nil {
				return output.Error(cmd, fmt.Sprintf("Error initiating deployment on new version %s", newVersionToDeploy), err)
			}
		} else {
			return output.Error(cmd, "Error deploying function", err)
		}
	}

	output.Success(cmd, fmt.Sprintf("Function %s version %s deployed successfully", functionId, versionId))
	output.Success(cmd, fmt.Sprintf("You can use 'nvcf function watch %s' to monitor the deployment", functionId))
	if !detached {
		return WaitForDeployment(cmd, client, functionId, versionId)
	}

	return nil
}

func createNewVersion(cmd *cobra.Command, client *api.Client, function nvcf.FunctionResponseFunction) (string, error) {
	containerEnv := mapResponseContainerEnvToNewContainerEnv(function.ContainerEnvironment)
	newVersion, err := client.Functions.Versions.New(cmd.Context(), function.ID, nvcf.FunctionVersionNewParams{
		Name:                 nvcf.String(function.Name),
		InferenceURL:         nvcf.String(function.InferenceURL),
		InferencePort:        nvcf.Int(function.InferencePort),
		ContainerImage:       nvcf.String(function.ContainerImage),
		ContainerArgs:        nvcf.String(function.ContainerArgs),
		ContainerEnvironment: nvcf.F(containerEnv),
		APIBodyFormat:        nvcf.F(nvcf.FunctionVersionNewParamsAPIBodyFormat(function.APIBodyFormat)),
		Description:          nvcf.F(function.Description),
		Tags:                 nvcf.F(function.Tags),
		FunctionType:         nvcf.F(nvcf.FunctionVersionNewParamsFunctionType(function.FunctionType)),
		Models:               nvcf.F(mapResponseModelsToNewModels(function.Models)),
		Health: nvcf.F(nvcf.FunctionVersionNewParamsHealth{
			Protocol:           nvcf.F(nvcf.FunctionVersionNewParamsHealthProtocol(function.Health.Protocol)),
			Port:               nvcf.F(function.Health.Port),
			Timeout:            nvcf.F(function.Health.Timeout),
			ExpectedStatusCode: nvcf.F(function.Health.ExpectedStatusCode),
			Uri:                nvcf.String(function.Health.Uri),
		}),
	})
	if err != nil {
		return "", output.Error(cmd, "Error creating new version", err)
	}

	return newVersion.Function.VersionID, nil
}

func mapResponseContainerEnvToNewContainerEnv(containerEnv []nvcf.FunctionResponseFunctionContainerEnvironment) []nvcf.FunctionVersionNewParamsContainerEnvironment {
	var newContainerEnv []nvcf.FunctionVersionNewParamsContainerEnvironment
	for _, env := range containerEnv {
		newContainerEnv = append(newContainerEnv, nvcf.FunctionVersionNewParamsContainerEnvironment{
			Key:   nvcf.F(env.Key),
			Value: nvcf.F(env.Value),
		})
	}
	return newContainerEnv
}

func mapResponseModelsToNewModels(models []nvcf.FunctionResponseFunctionModel) []nvcf.FunctionVersionNewParamsModel {
	var newModels []nvcf.FunctionVersionNewParamsModel
	for _, model := range models {
		newModels = append(newModels, nvcf.FunctionVersionNewParamsModel{
			Name:    nvcf.F(model.Name),
			Uri:     nvcf.F(model.Uri),
			Version: nvcf.F(model.Version),
		})
	}
	return newModels
}
