package function

import (
	"fmt"

	"github.com/brevdev/nvcf/api"
	"github.com/brevdev/nvcf/config"
	"github.com/brevdev/nvcf/output"
	"github.com/spf13/cobra"
)

func functionDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [function-id]",
		Short: "Delete a function. If you want to delete a specific version, use the --version-id flag.",
		Args:  cobra.ExactArgs(1),
		Run:   runFunctionDelete,
	}
	cmd.Flags().String("version-id", "", "The ID of the version")
	return cmd
}

func runFunctionDelete(cmd *cobra.Command, args []string) {
	client := api.NewClient(config.GetAPIKey())

	functionId := args[0]
	versionId, _ := cmd.Flags().GetString("version-id")

	if versionId == "" {
		output.Info(cmd, fmt.Sprintf("Deleting all versions of function %s", functionId))
		versions, err := client.Functions.Versions.List(cmd.Context(), functionId)
		if err != nil {
			output.Error(cmd, "Error listing function versions", err)
			return
		}
		for _, version := range versions.Functions {
			output.Info(cmd, fmt.Sprintf("Deleting version %s of function %s", version.VersionID, functionId))
			err := client.Functions.Versions.Delete(cmd.Context(), functionId, version.VersionID)
			if err != nil {
				output.Error(cmd, "Error deleting function version", err)
			}
		}
	} else {
		err := client.Functions.Versions.Delete(cmd.Context(), functionId, versionId)
		if err != nil {
			output.Error(cmd, "Error deleting function", err)
			return
		}
		output.Success(cmd, fmt.Sprintf("Function %s version %s deleted successfully", functionId, versionId))
	}
}
