// Package cmd provides command-line interface functionality for the NVCF CLI.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	// outputDir is the directory where generated documentation will be saved.
	outputDir string
)

// DocsCmd returns a cobra.Command for generating documentation for the NVCF CLI.
func DocsCmd() *cobra.Command {
	var genMarkdown, genReST bool
	cmd := &cobra.Command{
		Use:    "docs",
		Short:  "Generate documentation for the NVCF CLI",
		Long:   `Generate documentation in Markdown and/or reStructuredText formats for all NVCF CLI commands.`,
		Hidden: os.Getenv("NVCF_SHOW_DOCS_CMD") != "true", // Hide unless env var is set
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateDocs(cmd, genMarkdown, genReST)
		},
	}
	cmd.Flags().StringVarP(&outputDir, "output", "o", "./docs", "Output directory for generated documentation")
	cmd.Flags().BoolVar(&genMarkdown, "markdown", true, "Generate Markdown documentation")
	cmd.Flags().BoolVar(&genReST, "rst", false, "Generate reStructuredText documentation")
	return cmd
}

// generateDocs generates Markdown and/or reStructuredText documentation for the NVCF CLI.
func generateDocs(cmd *cobra.Command, genMarkdown, genReST bool) error {
	// Ensure the output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	// Get the root command
	root := cmd.Root()
	// Generate Markdown documentation
	if genMarkdown {
		if err := doc.GenMarkdownTree(root, outputDir); err != nil {
			return fmt.Errorf("failed to generate markdown documentation: %w", err)
		}
		fmt.Println("Markdown documentation generated")
	}
	// Generate reStructuredText documentation
	if genReST {
		emptyStr := func(s string) string { return "" }
		linkHandler := func(name, ref string) string {
			//sphinx-style refs.
			return fmt.Sprintf(":ref:`%s <%s>`", name, ref)
		}
		if err := doc.GenReSTTreeCustom(root, outputDir, emptyStr, linkHandler); err != nil {
			return fmt.Errorf("failed to generate reStructuredText documentation: %w", err)
		}
		fmt.Println("reStructuredText documentation generated")
	}
	fmt.Printf("Documentation generated in %s\n", outputDir)
	return nil
}
