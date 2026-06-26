package main

import (
	"fmt"
	"os"

	"github.com/ncobase/cli/commands"
)

func main() {
	rootCmd := commands.NewRootCmd()
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
