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
		Short: "A set of reusable components for Go applications",
	}

	// Add subcommands
	rootCmd.AddCommand(
		NewStartCommand(),
		NewPluginCommand(),
		NewDocsCommand(),
		NewVersionCommand(),
		NewCreateCommand(),
		NewInitCommand(),
		NewMigrateCommand(),
		NewSchemaCommand(),
	)

	return rootCmd
}
