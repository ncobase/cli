package migrate

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func newDiffCommand() *cobra.Command {
	var (
		dirPath    string
		to         string
		devURL     string
		configFile string
		env        string
	)

	cmd := &cobra.Command{
		Use:   "diff [name]",
		Short: "Compute the difference between the Ent schema and the migration directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			cArgs := []string{"migrate", "diff", name}

			if configFile != "" {
				cArgs = append(cArgs, "--config", configFile)
				if env != "" {
					cArgs = append(cArgs, "--env", env)
				}
			} else {
				// Default values
				if dirPath == "" {
					dirPath = "file://data/ent/migrate/migrations"
				}
				if to == "" {
					to = "ent://data/ent/schema"
				}
				if devURL == "" {
					devURL = "docker://mysql/8/ent" // Default to mysql
				}
				cArgs = append(cArgs, "--dir", dirPath, "--to", to, "--dev-url", devURL)
			}
			
			fmt.Printf("Running: atlas %v\n", cArgs)
			
			c := exec.Command("atlas", cArgs...)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin

			return c.Run()
		},
	}

	cmd.Flags().StringVar(&dirPath, "dir", "", "migration directory (default: file://data/ent/migrate/migrations)")
	cmd.Flags().StringVar(&to, "to", "", "target schema (default: ent://data/ent/schema)")
	cmd.Flags().StringVar(&devURL, "dev-url", "", "dev database URL (default: docker://mysql/8/ent)")
	cmd.Flags().StringVarP(&configFile, "config", "c", "", "atlas config file (default: atlas.hcl)")
	cmd.Flags().StringVarP(&env, "env", "e", "", "atlas env to use from config file")

	return cmd
}
