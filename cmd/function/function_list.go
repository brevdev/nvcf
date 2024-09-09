package function

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tmc/nvcf-go"
	"github.com/tmc/nvcf/api"
	"github.com/tmc/nvcf/config"
	"github.com/tmc/nvcf/output"
)

func functionListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all functions",
		Run:   runFunctionList,
	}
	cmd.Flags().StringSlice("visibility", []string{}, "Filter by visibility (authorized, private, public)")
	return cmd
}

func runFunctionList(cmd *cobra.Command, args []string) {
	client := api.NewClient(config.GetAPIKey())
	visibilityParams := parseVisibilityFlags(cmd)

	resp, err := client.Functions.List(cmd.Context(), nvcf.FunctionListParams{
		Visibility: nvcf.F(visibilityParams),
	})
	if err != nil {
		output.Error(cmd, "Error listing functions", err)
		return
	}
	output.Functions(cmd, resp.Functions)
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

// Implement other function subcommands (create, get, update, delete, version) here
