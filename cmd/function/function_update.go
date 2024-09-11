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

func functionUpdateCmd() *cobra.Command {
	var (
		gpu                   string
		instanceType          string
		minInstances          int64
		maxInstances          int64
		maxRequestConcurrency int64
	)
	cmd := &cobra.Command{
		Use:     "update <function-id>",
		Short:   "Update a function",
		Long:    "If a version-id is not provided, we look for versions that are actively deployed. If a single function is deployed, we update that version. If multiple functions are deployed, we prompt for the version-id to update.",
		Example: "nvcf function update fid --version-id vid --gpu A100 --instance-type g4dn.xlarge --min-instances 1 --max-instances 5 --max-request-concurrency 100",
		Args:    cobra.ExactArgs(1),
		Run:     runFunctionUpdate,
	}
	cmd.Flags().String("version-id", "", "The ID of the version")
	cmd.Flags().StringVar(&gpu, "gpu", "", "GPU name from the cluster")
	cmd.Flags().StringVar(&instanceType, "instance-type", "", "Instance type, based on GPU, assigned to a Worker")
	cmd.Flags().Int64Var(&minInstances, "min-instances", 0, "Minimum number of spot instances for the deployment")
	cmd.Flags().Int64Var(&maxInstances, "max-instances", 0, "Maximum number of spot instances for the deployment")
	cmd.Flags().Int64Var(&maxRequestConcurrency, "max-request-concurrency", 0, "Max request concurrency between 1 (default) and 1024")
	return cmd
}

func runFunctionUpdate(cmd *cobra.Command, args []string) {
	client := api.NewClient(config.GetAPIKey())
	functionID := args[0]
	versionID, _ := cmd.Flags().GetString("version-id")

	if versionID != "" {
		version, err := client.Functions.Versions.Get(cmd.Context(), functionID, versionID, nvcf.FunctionVersionGetParams{
			IncludeSecrets: nvcf.Bool(false),
		})
		if err != nil {
			output.Error(cmd, "Error getting function version", err)
			return
		}
		if version.Function.Status == "INACTIVE" {
			output.Error(cmd, "You can only update a deployed version. This version is inactive.", nil)
			return
		}
	}

	// Prepare update parameters
	var deploymentSpec nvcf.FunctionDeploymentFunctionVersionUpdateDeploymentParamsDeploymentSpecification

	if gpu, _ := cmd.Flags().GetString("gpu"); gpu != "" {
		deploymentSpec.GPU = nvcf.F(gpu)
	}
	if instanceType, _ := cmd.Flags().GetString("instance-type"); instanceType != "" {
		deploymentSpec.InstanceType = nvcf.F(instanceType)
	}
	if minInstances, _ := cmd.Flags().GetInt64("min-instances"); minInstances != 0 {
		deploymentSpec.MinInstances = nvcf.F(minInstances)
	}
	if maxInstances, _ := cmd.Flags().GetInt64("max-instances"); maxInstances != 0 {
		deploymentSpec.MaxInstances = nvcf.F(maxInstances)
	}
	if maxRequestConcurrency, _ := cmd.Flags().GetInt64("max-request-concurrency"); maxRequestConcurrency != 0 {
		deploymentSpec.MaxRequestConcurrency = nvcf.F(maxRequestConcurrency)
	}

	// Get all versions of the function
	versions, err := client.Functions.Versions.List(cmd.Context(), functionID)
	if err != nil {
		output.Error(cmd, "Error listing function versions", err)
		return
	}

	if len(versions.Functions) == 1 {
		versionID = versions.Functions[0].VersionID
		// make sure status is not inactive
		if versions.Functions[0].Status == "INACTIVE" {
			output.Error(cmd, "You can only update a deployed version. This version is inactive.", nil)
			return
		}
		updateParams := nvcf.FunctionDeploymentFunctionVersionUpdateDeploymentParams{
			DeploymentSpecifications: nvcf.F([]nvcf.FunctionDeploymentFunctionVersionUpdateDeploymentParamsDeploymentSpecification{deploymentSpec}),
		}

		_, err := client.FunctionDeployment.Functions.Versions.UpdateDeployment(cmd.Context(), functionID, versionID, updateParams)
		if err != nil {
			output.Error(cmd, "Error updating function deployment", err)
			return
		}

		output.Info(cmd, fmt.Sprintf("Successfully updated function deployment %s, version %s", functionID, versionID))
		return
	} else {
		output.Info(cmd, "Multiple deployed versions found. Please select a version to update:")
		for _, functions := range versions.Functions {
			if functions.Status != "INACTIVE" {
				output.Info(cmd, fmt.Sprintf("Version ID: %s || Status: %s", functions.VersionID, functions.Status))
			}
		}
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter version-id to update: ")
		versionID, _ = reader.ReadString('\n')
		versionID = strings.TrimSpace(versionID)
	}

	updateParams := nvcf.FunctionDeploymentFunctionVersionUpdateDeploymentParams{
		DeploymentSpecifications: nvcf.F([]nvcf.FunctionDeploymentFunctionVersionUpdateDeploymentParamsDeploymentSpecification{deploymentSpec}),
	}
	// Perform the update
	updatedFunction, err := client.FunctionDeployment.Functions.Versions.UpdateDeployment(cmd.Context(), functionID, versionID, updateParams)
	if err != nil {
		output.Error(cmd, "Error updating function deployment", err)
		return
	}

	output.Info(cmd, fmt.Sprintf("Successfully updated function deployment %s, version %s", functionID, versionID))
	output.SingleDeployment(cmd, *updatedFunction)
}
