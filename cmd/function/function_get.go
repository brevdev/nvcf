package function

import (
	"fmt"

	"github.com/spf13/cobra"
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

	cmd.MarkFlagRequired("function-id")
	cmd.MarkFlagRequired("version-id")
	return cmd
}

func runFunctionGet(cmd *cobra.Command, args []string) {
	functionID, _ := cmd.Flags().GetString("function-id")
	versionID, _ := cmd.Flags().GetString("version-id")
	output.Info(cmd, fmt.Sprintf("Getting function %s with version %s", functionID, versionID))
	client := api.NewClient(config.GetAPIKey())

	funcRes, err := client.FunctionManagement.Functions.Versions.Get(cmd.Context(), functionID, versionID)
	if err != nil {
		output.Error(cmd, "Error getting function", err)
		return
	}
	output.SingleFunction(cmd, funcRes.Function)
}
