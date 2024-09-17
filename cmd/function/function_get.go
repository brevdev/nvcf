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

func functionGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get <function-id>",
		Args:    cobra.ExactArgs(1),
		Short:   "Get details about a single function and its versions",
		Long:    "Get details about a single function and its versions or deployments. If a version-id is not provided and there are multiple versions associated with a function, we will look for all versions and prompt for a version-id.",
		Example: "nvcf function get fid --version-id vid --include-secrets",
		RunE:    runFunctionGet,
	}
	cmd.Flags().String("version-id", "", "The ID of the version")
	cmd.Flags().Bool("include-secrets", false, "Include secrets in the response")
	return cmd
}

func runFunctionGet(cmd *cobra.Command, args []string) error {
	client := api.NewClient(config.GetAPIKey())
	functionID := args[0]
	versionID, _ := cmd.Flags().GetString("version-id")
	includeSecrets, _ := cmd.Flags().GetBool("include-secrets")

	if versionID == "" {
		versions, err := client.Functions.Versions.List(cmd.Context(), functionID)
		if err != nil {
			output.Error(cmd, "Error listing function versions", err)
			return nil
		}

		if len(versions.Functions) == 1 {
			versionID = versions.Functions[0].VersionID
		} else {
			output.Info(cmd, "Multiple versions found. Please specify a version-id")
			for _, version := range versions.Functions {
				output.Info(cmd, fmt.Sprintf("Version ID: %s || Status: %s", version.VersionID, version.Status))
			}
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter version-id: ")
			versionID, _ = reader.ReadString('\n')
			versionID = strings.TrimSpace(versionID)
		}
	}

	query := nvcf.FunctionVersionGetParams{
		IncludeSecrets: nvcf.Bool(includeSecrets),
	}
	getFunctionResponse, err := client.Functions.Versions.Get(cmd.Context(), functionID, versionID, query)
	if err != nil {
		output.Error(cmd, "Error getting function", err)
		return nil
	}
	output.SingleFunction(cmd, getFunctionResponse.Function)
	return nil
}
