package output

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/tmc/nvcf-go"
)

func Error(cmd *cobra.Command, message string, err error) error {
	verbose, _ := cmd.Flags().GetBool("verbose")
	var formattedError error
	if verbose {
		formattedError = fmt.Errorf("%s: %v", message, err)
	} else {
		formattedError = fmt.Errorf("%s", message)
	}
	return formattedError
}

func Success(cmd *cobra.Command, message string) {
	if !isQuiet(cmd) {
		color.Green(message)
	}
}

func Info(cmd *cobra.Command, message string) {
	if !isQuiet(cmd) {
		color.Blue(message)
	}
}

func isJSON(cmd *cobra.Command) bool {
	json, _ := cmd.Flags().GetBool("json")
	return json
}

func isQuiet(cmd *cobra.Command) bool {
	quiet, _ := cmd.Flags().GetBool("quiet")
	return quiet
}

func printJSON(cmd *cobra.Command, data interface{}) {
	json, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		Error(cmd, "Error formatting JSON", err)
		return
	}
	fmt.Println(string(json))
}

func Prompt(message string, isSecret bool) string {
	fmt.Print(message)
	if isSecret {
		// Implement secure input for secrets
	}
	var input string
	fmt.Scanln(&input)
	return input
}

// Implement other output functions (Function, Deployment, InvocationResult, etc.) here
func Functions(cmd *cobra.Command, functions []nvcf.ListFunctionsResponseFunction) {
	if isJSON(cmd) {
		printJSON(cmd, functions)
	} else {
		printFunctionsTable(cmd, functions)
	}
}
func printFunctionsTable(cmd *cobra.Command, functions []nvcf.ListFunctionsResponseFunction) {
	table := tablewriter.NewWriter(cmd.OutOrStdout())
	table.SetHeader([]string{"Name", "Version ID", "Status"})
	table.SetBorder(false)
	for _, f := range functions {
		table.Append([]string{f.Name, f.VersionID, string(f.Status)})
	}
	table.Render()
}

func SingleFunction(cmd *cobra.Command, fn nvcf.FunctionResponseFunction) {
	if isJSON(cmd) {
		printJSON(cmd, fn)
	} else {
		printSingleFunctionTable(cmd, fn)
	}
}

func printSingleFunctionTable(cmd *cobra.Command, fn nvcf.FunctionResponseFunction) {
	table := tablewriter.NewWriter(cmd.OutOrStdout())
	table.SetHeader([]string{"Name", "Version ID", "Status"})
	table.SetBorder(false)
	table.Append([]string{fn.Name, fn.VersionID, string(fn.Status)})
	table.Render()
}

func Deployments(cmd *cobra.Command, deployments []nvcf.DeploymentResponse) {
	if isJSON(cmd) {
		printJSON(cmd, deployments)
	} else {
		printDeploymentsTable(cmd, deployments)
	}
}

func printDeploymentsTable(cmd *cobra.Command, deployments []nvcf.DeploymentResponse) {
	table := tablewriter.NewWriter(cmd.OutOrStdout())
	table.SetHeader([]string{"Function ID", "Function Version ID", "Status"})
	table.SetBorder(false)
	for _, deployment := range deployments {
		table.Append([]string{deployment.Deployment.FunctionID, deployment.Deployment.FunctionVersionID, string(deployment.Deployment.FunctionStatus)})
	}
	table.Render()
}

func SingleDeployment(cmd *cobra.Command, deployment nvcf.DeploymentResponse) {
	if isJSON(cmd) {
		printJSON(cmd, deployment)
	} else {
		printSingleDeploymentTable(cmd, deployment)
	}
}

func printSingleDeploymentTable(cmd *cobra.Command, deployment nvcf.DeploymentResponse) {
	table := tablewriter.NewWriter(cmd.OutOrStdout())
	table.SetHeader([]string{"Name", "Version ID", "Status"})
	table.SetBorder(false)
	table.Append([]string{deployment.Deployment.FunctionID, deployment.Deployment.FunctionVersionID, string(deployment.Deployment.FunctionStatus)})
	table.Render()
}

func GPUs(cmd *cobra.Command, clusterGroups []nvcf.ClusterGroupsResponseClusterGroup) {
	if isJSON(cmd) {
		printJSON(cmd, clusterGroups)
	} else {
		printGPUsTable(cmd, clusterGroups)
	}
}

func printGPUsTable(cmd *cobra.Command, clusterGroups []nvcf.ClusterGroupsResponseClusterGroup) {
	table := tablewriter.NewWriter(cmd.OutOrStdout())
	table.SetHeader([]string{"inst_backend", "inst_gpu_type", "inst_type"})
	table.SetBorder(false)

	for _, clusterGroup := range clusterGroups {
		for _, gpu := range clusterGroup.GPUs {
			for _, instanceType := range gpu.InstanceTypes {
				table.Append([]string{
					clusterGroup.Name, // inst_backend
					gpu.Name,          // inst_gpu_type
					instanceType.Name, // inst_type
				})
			}
		}
	}

	table.Render()
}
