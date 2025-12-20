package migrate

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func newNewCommand() *cobra.Command {
	var migrationsPath string

	cmd := &cobra.Command{
		Use:     "new [name]",
		Aliases: []string{"create"},
		Short:   "Create a new migration file",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			
			if migrationsPath == "" {
				migrationsPath = "file://data/ent/migrate/migrations"
			}

			// atlas migrate new name --dir ...
			cArgs := []string{"migrate", "new", name, "--dir", migrationsPath}
			
			fmt.Printf("Running: atlas %v\n", cArgs)
			
			c := exec.Command("atlas", cArgs...)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin

			return c.Run()
		},
	}

	cmd.Flags().StringVarP(&migrationsPath, "dir", "d", "", "migrations directory path (default: file://data/ent/migrate/migrations)")
	return cmd
}
