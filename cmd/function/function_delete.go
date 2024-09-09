package function

import (
	"github.com/spf13/cobra"
	"github.com/tmc/nvcf/api"
	"github.com/tmc/nvcf/config"
	"github.com/tmc/nvcf/output"
)

func functionDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a function",
		Run:   runFunctionDelete,
	}
	cmd.Flags().String("function-id", "", "The ID of the function")
	cmd.Flags().String("version-id", "", "The ID of the version")

	cmd.MarkFlagRequired("function-id")
	cmd.MarkFlagRequired("version-id")
	return cmd
}

func runFunctionDelete(cmd *cobra.Command, args []string) {
	functionID, _ := cmd.Flags().GetString("function-id")
	versionID, _ := cmd.Flags().GetString("version-id")

	client := api.NewClient(config.GetAPIKey())
	err := client.Functions.Versions.Delete(cmd.Context(), functionID, versionID)
	if err != nil {
		output.Error(cmd, "Error deleting function", err)
		return
	}
	output.Success(cmd, "Function deleted successfully")
}
