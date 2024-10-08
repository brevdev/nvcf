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
)

func functionDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <function-id>",
		Short:   "Delete a function. If you want to delete a specific version, use the --version-id flag.",
		Long:    "Delete a function. If there is only 1 version, we will delete the function. If there are multiple versions, we will prompt you to specify which version to delete. The --all flag will delete all versions of the function. Deleting a deployed function will change a function status to INACTIVE and using the --force flag will delete the function immediately.",
		Example: "nvcf function delete fid --version-id vid",
		Args:    cobra.ExactArgs(1),
		RunE:    runFunctionDelete,
	}
	cmd.Flags().String("version-id", "", "The ID of the version")
	cmd.Flags().BoolP("all", "a", false, "Delete all versions of the function")
	cmd.Flags().BoolP("force", "f", false, "Forcefully delete a deployed function")
	return cmd
}

func runFunctionDelete(cmd *cobra.Command, args []string) error {
	client := api.NewClient(config.GetAPIKey())

	functionId := args[0]
	versionId, _ := cmd.Flags().GetString("version-id")
	all, _ := cmd.Flags().GetBool("all")

	if versionId == "" {
		versions, err := client.Functions.Versions.List(cmd.Context(), functionId)
		if err != nil {
			return output.Error(cmd, "Error listing function versions", err)
		}
		if len(versions.Functions) == 1 {
			versionId = versions.Functions[0].VersionID
			output.Info(cmd, "Deleting function")
			err := client.Functions.Versions.Delete(cmd.Context(), functionId, versionId)
			if err != nil {
				return output.Error(cmd, "Error deleting function", err)
			}
			output.Success(cmd, "Function deleted successfully")
			return nil
		} else {
			if all {
				for _, version := range versions.Functions {
					output.Info(cmd, fmt.Sprintf("Deleting version %s of function %s", version.VersionID, functionId))
					err := client.Functions.Versions.Delete(cmd.Context(), functionId, version.VersionID)
					if err != nil {
						return output.Error(cmd, "Error deleting function version", err)
					}
				}
				output.Success(cmd, "All versions of function deleted successfully")
				return nil
			} else {
				output.Info(cmd, "Multiple versions found. Please specify a version-id")
				for _, version := range versions.Functions {
					output.Info(cmd, fmt.Sprintf("Version ID: %s || Status: %s", version.VersionID, version.Status))
				}
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Enter version-id: ")
				versionId, _ = reader.ReadString('\n')
				versionId = strings.TrimSpace(versionId)
				err := client.Functions.Versions.Delete(cmd.Context(), functionId, versionId)
				if err != nil {
					return output.Error(cmd, "Error deleting function version", err)
				}
				output.Success(cmd, "Function deleted successfully")
			}
		}
	} else {
		output.Info(cmd, fmt.Sprintf("Deleting function %s version %s", functionId, versionId))
		err := client.Functions.Versions.Delete(cmd.Context(), functionId, versionId)
		if err != nil {
			return output.Error(cmd, "Error deleting function", err)
		}
		output.Success(cmd, "Function deleted successfully")
	}
	return nil
}
