package function

import (
	"fmt"

	"github.com/brevdev/nvcf/api"
	"github.com/brevdev/nvcf/config"
	"github.com/brevdev/nvcf/output"
	"github.com/spf13/cobra"
	"github.com/tmc/nvcf-go"
)

func functionGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get [identifier]",
		Args:    cobra.MaximumNArgs(1),
		Short:   "Get details about a single function and its versions",
		Long:    "Get details about a single function and its versions or deployments. The identifier can be a function name, function ID, or version ID. If no identifier is provided, all functions will be listed.",
		Example: "nvcf function get myFunction\nnvcf function get fid123\nnvcf function get --name myFunction\nnvcf function get --function-id fid123 --version-id vid456",
		RunE:    runFunctionGet,
	}
	cmd.Flags().String("name", "", "Filter by function name")
	cmd.Flags().String("function-id", "", "Filter by function ID")
	cmd.Flags().String("version-id", "", "Filter by version ID")
	cmd.Flags().Bool("include-secrets", false, "Include secrets in the response")
	return cmd
}

func runFunctionGet(cmd *cobra.Command, args []string) error {
	client := api.NewClient(config.GetAPIKey())
	includeSecrets, _ := cmd.Flags().GetBool("include-secrets")

	name, _ := cmd.Flags().GetString("name")
	functionID, _ := cmd.Flags().GetString("function-id")
	versionID, _ := cmd.Flags().GetString("version-id")

	identifier := ""
	if len(args) > 0 {
		identifier = args[0]
	}

	if identifier == "" && name == "" && functionID == "" && versionID == "" {
		return listAllFunctions(cmd, client)
	}

	functions, err := client.Functions.List(cmd.Context(), nvcf.FunctionListParams{})
	if err != nil {
		return output.Error(cmd, "Error listing functions", err)
	}

	matchedFunctions := []nvcf.ListFunctionsResponseFunction{}
	for _, fn := range functions.Functions {
		if matchesIdentifier(fn, identifier, name, functionID, versionID) {
			query := nvcf.FunctionVersionGetParams{
				IncludeSecrets: nvcf.Bool(includeSecrets),
			}
			_, err := client.Functions.Versions.Get(cmd.Context(), fn.ID, fn.VersionID, query)
			if err != nil {
				return output.Error(cmd, fmt.Sprintf("Error getting function %s", fn.ID), err)
			}
			matchedFunctions = append(matchedFunctions, fn)
		}
	}

	if len(matchedFunctions) == 0 {
		return output.Error(cmd, "No matching functions found", nil)
	}

	output.Functions(cmd, matchedFunctions)
	return nil
}

func matchesIdentifier(fn nvcf.ListFunctionsResponseFunction, identifier, name, functionID, versionID string) bool {
	return (identifier != "" && (fn.Name == identifier || fn.ID == identifier || fn.VersionID == identifier)) ||
		(name != "" && fn.Name == name) ||
		(functionID != "" && fn.ID == functionID) ||
		(versionID != "" && fn.VersionID == versionID)
}

func listAllFunctions(cmd *cobra.Command, client *api.Client) error {
	functions, err := client.Functions.List(cmd.Context(), nvcf.FunctionListParams{})
	if err != nil {
		return output.Error(cmd, "Error listing functions", err)
	}
	output.Functions(cmd, functions.Functions)
	return nil
}
