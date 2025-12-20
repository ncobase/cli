package migrate

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func newStatusCommand() *cobra.Command {
	var (
		migrationsPath string
		databaseURL    string
		configFile     string
		env            string
	)

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Get information about the migration status of the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			cArgs := []string{"migrate", "status"}

			if configFile != "" {
				cArgs = append(cArgs, "--config", configFile)
				if env != "" {
					cArgs = append(cArgs, "--env", env)
				}
			} else {
				if migrationsPath == "" {
					migrationsPath = "file://data/ent/migrate/migrations"
				}
				if databaseURL == "" {
					return fmt.Errorf("database url is required (use -u or --url)")
				}
				cArgs = append(cArgs, "--dir", migrationsPath, "--url", databaseURL)
			}

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
	cmd.Flags().StringVarP(&configFile, "config", "c", "", "atlas config file (default: atlas.hcl)")
	cmd.Flags().StringVarP(&env, "env", "e", "", "atlas env to use from config file")

	return cmd
}
