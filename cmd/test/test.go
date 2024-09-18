package test

import (
	"fmt"

	"github.com/spf13/cobra"
)

func TestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "A simple test command",
		Long:  `This command simply returns "Hello, World!" when executed.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello, World!")
		},
	}
}
