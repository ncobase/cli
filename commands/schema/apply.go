package schema

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func newApplyCommand() *cobra.Command {
	var (
		databaseURL string
		to          string
		devURL      string
		configFile  string
		env         string
		dryRun      bool
		autoApprove bool
	)

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply a schema to the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			cArgs := []string{"schema", "apply"}

			if configFile != "" {
				cArgs = append(cArgs, "--config", configFile)
				if env != "" {
					cArgs = append(cArgs, "--env", env)
				}
			} else {
				if databaseURL == "" {
					return fmt.Errorf("database url is required (use -u or --url)")
				}
				if to == "" {
					// Default to local schema file or directory, but user should specify
					return fmt.Errorf("target schema is required (use --to)")
				}
				cArgs = append(cArgs, "--url", databaseURL, "--to", to)
				
				if devURL != "" {
					cArgs = append(cArgs, "--dev-url", devURL)
				}
			}

			if dryRun {
				cArgs = append(cArgs, "--dry-run")
			}
			if autoApprove {
				cArgs = append(cArgs, "--auto-approve")
			}

			fmt.Printf("Running: atlas %v\n", cArgs)
			
			c := exec.Command("atlas", cArgs...)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin

			return c.Run()
		},
	}

	cmd.Flags().StringVarP(&databaseURL, "url", "u", "", "database url (e.g. mysql://user:pass@localhost:3306/db)")
	cmd.Flags().StringVar(&to, "to", "", "target schema (e.g. file://schema.hcl)")
	cmd.Flags().StringVar(&devURL, "dev-url", "", "dev database URL (e.g. docker://mysql/8/dev)")
	cmd.Flags().StringVarP(&configFile, "config", "c", "", "atlas config file (default: atlas.hcl)")
	cmd.Flags().StringVarP(&env, "env", "e", "", "atlas env to use from config file")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "dry run")
	cmd.Flags().BoolVar(&autoApprove, "auto-approve", false, "auto approve")

	return cmd
}
