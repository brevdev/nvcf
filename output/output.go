package output

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/tmc/nvcf-go"
)

func Error(cmd *cobra.Command, message string, err error) {
	if !isQuiet(cmd) {
		color.Red(fmt.Sprintf("%s: %v", message, err))
	}
}

func Success(cmd *cobra.Command, message string) {
	if !isQuiet(cmd) {
		color.Green(message)
	}
}

func Info(cmd *cobra.Command, message string) {
	if !isQuiet(cmd) {
		color.Blue(message)
	}
}

func isJSON(cmd *cobra.Command) bool {
	json, _ := cmd.Flags().GetBool("json")
	return json
}

func isQuiet(cmd *cobra.Command) bool {
	quiet, _ := cmd.Flags().GetBool("quiet")
	return quiet
}

func printJSON(cmd *cobra.Command, data interface{}) {
	json, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		Error(cmd, "Error formatting JSON", err)
		return
	}
	fmt.Println(string(json))
}

func Prompt(message string, isSecret bool) string {
	fmt.Print(message)
	if isSecret {
		// Implement secure input for secrets
	}
	var input string
	fmt.Scanln(&input)
	return input
}

// Implement other output functions (Function, Deployment, InvocationResult, etc.) here
func Functions(cmd *cobra.Command, functions []nvcf.ListFunctionsResponseFunction) {
	if isJSON(cmd) {
		printJSON(cmd, functions)
	} else {
		printFunctionsTable(cmd, functions)
	}
}
func printFunctionsTable(cmd *cobra.Command, functions []nvcf.ListFunctionsResponseFunction) {
	table := tablewriter.NewWriter(cmd.OutOrStdout())
	table.SetHeader([]string{"Name", "ID", "Status"})
	table.SetBorder(false)
	for _, f := range functions {
		table.Append([]string{f.Name, f.ID, string(f.Status)})
	}
	table.Render()
}

func MultiFunction(cmd *cobra.Command, fn nvcf.ListFunctionsResponseFunction) {
	if isJSON(cmd) {
		printJSON(cmd, fn)
	} else {
		printMultiFunctionTable(cmd, fn)
	}
}
func printMultiFunctionTable(cmd *cobra.Command, fn nvcf.ListFunctionsResponseFunction) {
	table := tablewriter.NewWriter(cmd.OutOrStdout())
	table.SetHeader([]string{"Name", "ID", "Status"})
	table.SetBorder(false)
	table.Append([]string{fn.Name, fn.ID, string(fn.Status)})
	table.Render()
}

func SingleFunction(cmd *cobra.Command, fn nvcf.FunctionResponseFunction) {
	if isJSON(cmd) {
		printJSON(cmd, fn)
	} else {
		printSingleFunctionTable(cmd, fn)
	}
}

func printSingleFunctionTable(cmd *cobra.Command, fn nvcf.FunctionResponseFunction) {
	table := tablewriter.NewWriter(cmd.OutOrStdout())
	table.SetHeader([]string{"Name", "ID", "Status"})
	table.SetBorder(false)
	table.Append([]string{fn.Name, fn.ID, string(fn.Status)})
	table.Render()
}
