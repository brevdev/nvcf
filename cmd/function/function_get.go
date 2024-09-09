package function

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tmc/nvcf-go"
	"github.com/tmc/nvcf/api"
	"github.com/tmc/nvcf/config"
	"github.com/tmc/nvcf/output"
)

func functionGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [function-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Get details about a single function. If you want to get a specific version, use the --version-id flag.",
		Run:   runFunctionGet,
	}
	cmd.Flags().String("version-id", "", "The ID of the version")
	cmd.Flags().Bool("include-secrets", false, "Include secrets in the response")
	return cmd
}

func runFunctionGet(cmd *cobra.Command, args []string) {
	client := api.NewClient(config.GetAPIKey())
	functionID := args[0]
	versionID, _ := cmd.Flags().GetString("version-id")
	includeSecrets, _ := cmd.Flags().GetBool("include-secrets")

	if versionID == "" {
		output.Info(cmd, fmt.Sprintf("Getting all versions of function %s", functionID))
		versions, err := client.Functions.Versions.List(cmd.Context(), functionID)
		if err != nil {
			output.Error(cmd, "Error listing function versions", err)
			return
		}
		for _, version := range versions.Functions {
			output.MultiFunction(cmd, version)
		}
	} else {
		output.Info(cmd, fmt.Sprintf("Getting version %s of function %s", versionID, functionID))
		query := nvcf.FunctionVersionGetParams{
			IncludeSecrets: nvcf.Bool(includeSecrets),
		}
		getFunctionResponse, err := client.Functions.Versions.Get(cmd.Context(), functionID, versionID, query)
		if err != nil {
			output.Error(cmd, "Error getting function", err)
			return
		}
		output.SingleFunction(cmd, getFunctionResponse.Function)
	}

}
