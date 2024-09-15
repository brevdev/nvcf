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
		Short:   "Update a deployed function",
		Long:    "If a version-id is not provided, we look for versions that are actively deployed. If a single function is deployed, we update that version. If multiple functions are deployed, we prompt for the version-id to update.",
		Example: "nvcf function update fid --version-id vid --gpu A100 --instance-type g5.4xlarge --min-instances 1 --max-instances 5 --max-request-concurrency 100",
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

	versions, err := client.Functions.Versions.List(cmd.Context(), functionID)
	if err != nil {
		output.Error(cmd, "Error listing function versions", err)
		return
	}

	var fnDeployment *nvcf.DeploymentResponse
	// Select the version to update
	if versionID == "" {
		vid := selectVersionToUpdate(cmd, versions.Functions)
		fnDeployment, err = client.FunctionDeployment.Functions.Versions.GetDeployment(cmd.Context(), functionID, vid)
		if err != nil {
			output.Error(cmd, "Error getting function version", err)
			return
		}
	} else {
		statusCheck, err := client.Functions.Versions.Get(cmd.Context(), functionID, versionID, nvcf.FunctionVersionGetParams{
			IncludeSecrets: nvcf.Bool(false),
		})
		if err != nil {
			output.Error(cmd, "Error getting function version", err)
			return
		}
		if statusCheck.Function.Status == "INACTIVE" {
			output.Error(cmd, "You can only update a deployed version. This version is inactive.", nil)
			return
		}
		fnDeployment, err = client.FunctionDeployment.Functions.Versions.GetDeployment(cmd.Context(), functionID, versionID)
		if err != nil {
			output.Error(cmd, "Error getting function version", err)
			return
		}
	}

	var targetDeployment nvcf.DeploymentResponseDeploymentDeploymentSpecification
	// there can be multiple deployments for a version - prompt and check similar to how we do the version check
	if len(fnDeployment.Deployment.DeploymentSpecifications) > 1 {
		output.Info(cmd, "Multiple deployment specifications found. Please select one to update:")
		for i, spec := range fnDeployment.Deployment.DeploymentSpecifications {
			output.Info(cmd, fmt.Sprintf("[%d] GPU: %s, Instance Type: %s, Min Instances: %d, Max Instances: %d, Max Request Concurrency: %d",
				i+1, spec.GPU, spec.InstanceType, spec.MinInstances, spec.MaxInstances, spec.MaxRequestConcurrency))
		}

		var selectedIndex int
		for {
			fmt.Print("Enter the number of the deployment specification to update: ")
			_, err := fmt.Scanf("%d", &selectedIndex)
			if err == nil && selectedIndex > 0 && selectedIndex <= len(fnDeployment.Deployment.DeploymentSpecifications) {
				break
			}
			output.Info(cmd, "Invalid selection. Please try again.")
		}

		targetDeployment = fnDeployment.Deployment.DeploymentSpecifications[selectedIndex-1]
	} else {
		targetDeployment = fnDeployment.Deployment.DeploymentSpecifications[0]
	}

	// build deployment spec with values from the targetDeployment
	deploymentSpec := nvcf.FunctionDeploymentFunctionVersionUpdateDeploymentParamsDeploymentSpecification{
		GPU:                   nvcf.String(targetDeployment.GPU),
		InstanceType:          nvcf.String(targetDeployment.InstanceType),
		MinInstances:          nvcf.Int(targetDeployment.MinInstances),
		MaxInstances:          nvcf.Int(targetDeployment.MaxInstances),
		MaxRequestConcurrency: nvcf.Int(targetDeployment.MaxRequestConcurrency),
	}

	if gpu, _ := cmd.Flags().GetString("gpu"); gpu != "" {
		deploymentSpec.GPU = nvcf.String(gpu)
	}
	if instanceType, _ := cmd.Flags().GetString("instance-type"); instanceType != "" {
		deploymentSpec.InstanceType = nvcf.String(instanceType)
	}
	if minInstances, _ := cmd.Flags().GetInt64("min-instances"); minInstances != 0 {
		deploymentSpec.MinInstances = nvcf.Int(minInstances)
	}
	if maxInstances, _ := cmd.Flags().GetInt64("max-instances"); maxInstances != 0 {
		deploymentSpec.MaxInstances = nvcf.Int(maxInstances)
	}
	if maxRequestConcurrency, _ := cmd.Flags().GetInt64("max-request-concurrency"); maxRequestConcurrency != 0 {
		deploymentSpec.MaxRequestConcurrency = nvcf.Int(maxRequestConcurrency)
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

func selectVersionToUpdate(cmd *cobra.Command, versions []nvcf.ListFunctionsResponseFunction) string {
	if len(versions) == 1 {
		if versions[0].Status == "INACTIVE" {
			output.Error(cmd, "You can only update a deployed version. This version is inactive.", nil)
			os.Exit(1)
		}
		return versions[0].VersionID
	}

	output.Info(cmd, "Multiple deployed versions found. Please select a version to update:")
	for _, function := range versions {
		if function.Status != "INACTIVE" {
			output.Info(cmd, fmt.Sprintf("Version ID: %s || Status: %s", function.VersionID, function.Status))
		}
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter version-id to update: ")
	versionID, _ := reader.ReadString('\n')
	return strings.TrimSpace(versionID)
}
