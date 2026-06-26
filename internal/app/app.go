package app

import (
	"github.com/ncobase/cli/commands"
)

// Execute runs the CLI root command.
func Execute() error {
	rootCmd := commands.NewRootCmd()
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	return rootCmd.Execute()
}
