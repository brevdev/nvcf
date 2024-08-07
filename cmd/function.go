package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tmc/nvcf/api"
	"github.com/tmc/nvcf/output"
)

func FunctionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "function",
		Short: "Manage NVIDIA Cloud Functions",
		Long:  `Create, list, update, and delete NVIDIA Cloud Functions.`,
	}

	cmd.AddCommand(functionListCmd())
	cmd.AddCommand(functionCreateCmd())
	cmd.AddCommand(functionGetCmd())
	cmd.AddCommand(functionUpdateCmd())
	cmd.AddCommand(functionDeleteCmd())
	cmd.AddCommand(functionVersionCmd())

	return cmd
}

func functionListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all functions",
		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewClient()
			limit, _ := cmd.Flags().GetInt("limit")
			visibility, _ := cmd.Flags().GetString("visibility")
			functions, err := client.ListFunctions(limit, visibility)
			if err != nil {
				output.Error(cmd, "Error listing functions", err)
				return
			}
			output.Functions(cmd, functions)
		},
	}

	cmd.Flags().Int("limit", 0, "Maximum number of functions to list")
	cmd.Flags().String("visibility", "", "Filter by visibility (authorized, private, public)")

	return cmd
}

// Implement other function subcommands (create, get, update, delete, version) here
