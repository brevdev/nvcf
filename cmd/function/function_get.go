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
		Use:   "get",
		Short: "Get details about a function",
		Run:   runFunctionGet,
	}
	cmd.Flags().String("function-id", "", "The ID of the function")
	cmd.Flags().String("version-id", "", "The ID of the version")
	cmd.Flags().Bool("include-secrets", false, "Include secrets in the response")

	cmd.MarkFlagRequired("function-id")
	cmd.MarkFlagRequired("version-id")
	return cmd
}

func runFunctionGet(cmd *cobra.Command, args []string) {
	functionID, _ := cmd.Flags().GetString("function-id")
	versionID, _ := cmd.Flags().GetString("version-id")
	includeSecrets, _ := cmd.Flags().GetBool("include-secrets")
	output.Info(cmd, fmt.Sprintf("Getting function %s with version %s", functionID, versionID))
	client := api.NewClient(config.GetAPIKey())

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
