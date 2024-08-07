package output

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
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

func Functions(cmd *cobra.Command, functions interface{}) {
	if isJSON(cmd) {
		printJSON(cmd, functions)
	} else {
		// Implement table output for functions
	}
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

// Implement other output functions (Function, Deployment, InvocationResult, etc.) here
