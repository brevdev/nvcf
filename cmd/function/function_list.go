package function

import (
	"fmt"

	"github.com/brevdev/nvcf/api"
	"github.com/brevdev/nvcf/config"
	"github.com/brevdev/nvcf/output"
	"github.com/spf13/cobra"
	"github.com/tmc/nvcf-go"
)

func functionListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "l"},
		Short:   "List all functions. Use flags to filter by visibility and status.",
		Long:    "List all functions. Use the --visibility flag to filter by visibility and --status flag to filter by status. Note that deployments are functions that are NOT inactive.",
		Example: "nvcf function list --visibility authorized --status ACTIVE,DEPLOYING",
		Run:     runFunctionList,
	}
	cmd.Flags().StringSlice("visibility", []string{"private"}, "Filter by visibility (authorized, private, public). Defaults to private.")
	cmd.Flags().StringSlice("status", []string{"ACTIVE", "DEPLOYING", "ERROR", "INACTIVE", "DELETED"}, "Filter by status (ACTIVE, DEPLOYING, ERROR, INACTIVE, DELETED). Defaults to all.")
	return cmd
}

func runFunctionList(cmd *cobra.Command, args []string) {
	client := api.NewClient(config.GetAPIKey())
	visibilityParams := parseVisibilityFlags(cmd)
	statusParams := parseStatusFlags(cmd)

	// Get all functions matching the visibility flags
	functions, err := client.Functions.List(cmd.Context(), nvcf.FunctionListParams{
		Visibility: nvcf.F(visibilityParams),
	})
	if err != nil {
		output.Error(cmd, "Error listing functions", err)
		return
	}

	// Filter functions by status params
	filteredFunctions := make([]nvcf.ListFunctionsResponseFunction, 0)
	for _, function := range functions.Functions {
		if containsStatus(statusParams, function.Status) {
			filteredFunctions = append(filteredFunctions, function)
		}
	}

	output.Functions(cmd, filteredFunctions)
}

func parseVisibilityFlags(cmd *cobra.Command) []nvcf.FunctionListParamsVisibility {
	visibilityFlags, _ := cmd.Flags().GetStringSlice("visibility")
	var visibilityParams []nvcf.FunctionListParamsVisibility
	for _, v := range visibilityFlags {
		param := nvcf.FunctionListParamsVisibility(v)
		if param.IsKnown() {
			visibilityParams = append(visibilityParams, param)
		} else {
			output.Error(cmd, fmt.Sprintf("Invalid visibility: '%s'", v), nil)
			return nil
		}
	}
	return visibilityParams
}

func parseStatusFlags(cmd *cobra.Command) []nvcf.DeploymentResponseDeploymentFunctionStatus {
	statusFlags, _ := cmd.Flags().GetStringSlice("status")
	var statusParams []nvcf.DeploymentResponseDeploymentFunctionStatus
	for _, v := range statusFlags {
		param := nvcf.DeploymentResponseDeploymentFunctionStatus(v)
		if param.IsKnown() {
			statusParams = append(statusParams, param)
		} else {
			output.Error(cmd, fmt.Sprintf("Invalid status: '%s'", v), nil)
			return nil
		}
	}
	return statusParams
}

func containsStatus(statuses []nvcf.DeploymentResponseDeploymentFunctionStatus, status nvcf.ListFunctionsResponseFunctionsStatus) bool {
	for _, s := range statuses {
		if string(s) == string(status) {
			return true
		}
	}
	return false
}
