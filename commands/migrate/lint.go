package migrate

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func newLintCommand() *cobra.Command {
	var (
		migrationsPath string
		devURL         string
		configFile     string
		env            string
		latest         int
	)

	cmd := &cobra.Command{
		Use:   "lint",
		Short: "Lint the migration directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			cArgs := []string{"migrate", "lint"}

			if configFile != "" {
				cArgs = append(cArgs, "--config", configFile)
				if env != "" {
					cArgs = append(cArgs, "--env", env)
				}
			} else {
				if migrationsPath == "" {
					migrationsPath = "file://data/ent/migrate/migrations"
				}
				if devURL == "" {
					devURL = "docker://mysql/8/ent" // Default to mysql
				}
				cArgs = append(cArgs, "--dir", migrationsPath, "--dev-url", devURL)
			}

			if latest > 0 {
				cArgs = append(cArgs, "--latest", fmt.Sprintf("%d", latest))
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
	cmd.Flags().StringVar(&devURL, "dev-url", "", "dev database URL (default: docker://mysql/8/ent)")
	cmd.Flags().StringVarP(&configFile, "config", "c", "", "atlas config file (default: atlas.hcl)")
	cmd.Flags().StringVarP(&env, "env", "e", "", "atlas env to use from config file")
	cmd.Flags().IntVar(&latest, "latest", 0, "lint only the latest N files")

	return cmd
}
