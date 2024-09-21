package test

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

func TestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "A simple test command",
		Long:  `This command simply returns "Hello, World!" when executed.`,
		Run: func(cmd *cobra.Command, args []string) {
			s := spinner.New(spinner.CharSets[4], 100*time.Millisecond)
			s.Suffix = " Hello, World!"
			s.Start()
			time.Sleep(4 * time.Second)
			s.Stop()
			fmt.Println("Finished!")
		},
	}
}
