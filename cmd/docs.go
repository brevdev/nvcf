package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	outputDir string
)

func DocsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "docs",
		Short:  "Generate documentation for the NVCF CLI",
		Long:   `Generate documentation in Markdown format for all NVCF CLI commands.`,
		Hidden: os.Getenv("NVCF_SHOW_DOCS_CMD") != "true", // Hide unless env var is set
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateDocs(cmd)
		},
	}

	cmd.Flags().StringVarP(&outputDir, "output", "o", "./docs", "Output directory for generated documentation")

	return cmd
}

func generateDocs(cmd *cobra.Command) error {
	// Ensure the output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Get the root command
	root := cmd.Root()

	// Generate Markdown documentation
	err := doc.GenMarkdownTree(root, outputDir)
	if err != nil {
		return fmt.Errorf("failed to generate markdown documentation: %w", err)
	}

	fmt.Printf("Documentation generated in %s\n", outputDir)
	return nil
}
