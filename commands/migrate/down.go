package migrate

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func newDownCommand() *cobra.Command {
	var (
		migrationsPath string
		databaseURL    string
	)

	cmd := &cobra.Command{
		Use:   "down",
		Short: "Rollback the last migration",
		RunE: func(cmd *cobra.Command, args []string) error {
			if migrationsPath == "" {
				migrationsPath = "file://data/ent/migrate/migrations"
			}

			if databaseURL == "" {
				return fmt.Errorf("database url is required (use -u or --url)")
			}

			// atlas migrate down --dir ... --url ...
			cArgs := []string{"migrate", "down", "--dir", migrationsPath, "--url", databaseURL}

			fmt.Printf("Running: atlas %v\n", cArgs)

			c := exec.Command("atlas", cArgs...)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin

			return c.Run()
		},
	}

	cmd.Flags().StringVarP(&migrationsPath, "dir", "d", "", "migrations directory path (default: file://data/ent/migrate/migrations)")
	cmd.Flags().StringVarP(&databaseURL, "url", "u", "", "database url (e.g. mysql://user:pass@localhost:3306/db)")
	return cmd
}
