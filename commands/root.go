package commands

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// NewRootCmd creates the root command
func NewRootCmd() *cobra.Command {
	// Define root command
	rootCmd := &cobra.Command{
		Use:   filepath.Base(os.Args[0]),
		Short: "Ncobase CLI - Scaffold Go applications with ncore framework",
		Long: `Ncobase CLI - A powerful scaffolding tool for the ncore framework

Core Commands:
  init     Initialize a new standalone application
  create   Create extensions within an existing project (core/business/plugin)

Database Utilities (Optional):
  migrate  Run database migrations (requires Atlas)
  schema   Inspect and manage database schemas (requires Atlas)

Use "nco [command] --help" for more information about a command.`,
	}

	// Add subcommands
	rootCmd.AddCommand(
		NewVersionCommand(),
		NewInitCommand(),
		NewCreateCommand(),
		NewMigrateCommand(),
		NewSchemaCommand(),
	)

	return rootCmd
}
