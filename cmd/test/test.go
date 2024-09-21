package test

import (
	"time"

	"github.com/brevdev/nvcf/output"
	"github.com/spf13/cobra"
)

func TestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "A simple test command",
		Long:  `This command simply returns "Hello, World!" when executed.`,
		Run: func(cmd *cobra.Command, args []string) {
			output.PrintASCIIArt(cmd)
			time.Sleep(1 * time.Second)
			s := output.NewSpinner(" Hello, World!")
			output.StartSpinner(s)
			time.Sleep(4 * time.Second)
			output.StopSpinner(s)
			output.Success(cmd, "Finished!")
		},
	}
}