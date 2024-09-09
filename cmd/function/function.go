package function

import (
	"fmt"
	"os"

	"github.com/brevdev/nvcf/api"
	"github.com/brevdev/nvcf/config"
	"github.com/brevdev/nvcf/output"
	"github.com/spf13/cobra"
	"github.com/tmc/nvcf-go"
)

func FunctionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "functions",
		Aliases: []string{"fn", "fns", "cf"},
		Short:   "Manage NVIDIA Cloud Functions",
		Long:    `Create, list, update, call, deploy, and delete NVIDIA Cloud Functions.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			config.Init()
			if cmd.Name() != "auth" && !config.IsAuthenticated() {
				fmt.Println("You are not authenticated. Please run 'nvcf auth login' first.")
				os.Exit(1)
			}
		},
	}

	cmd.AddCommand(functionListCmd())
	cmd.AddCommand(functionCreateCmd())
	// cmd.AddCommand(functionGetCmd())
	// cmd.AddCommand(functionUpdateCmd())
	// cmd.AddCommand(functionDeleteCmd())
	// cmd.AddCommand(functionVersionCmd())

	return cmd
}

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
