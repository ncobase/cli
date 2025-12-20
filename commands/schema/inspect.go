package schema

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func newInspectCommand() *cobra.Command {
	var (
		databaseURL string
		configFile  string
		env         string
		format      string
	)

	cmd := &cobra.Command{
		Use:   "inspect",
		Short: "Inspect a database and return its schema in HCL",
		RunE: func(cmd *cobra.Command, args []string) error {
			cArgs := []string{"schema", "inspect"}

			if configFile != "" {
				cArgs = append(cArgs, "--config", configFile)
				if env != "" {
					cArgs = append(cArgs, "--env", env)
				}
			} else {
				if databaseURL == "" {
					return fmt.Errorf("database url is required (use -u or --url)")
				}
				cArgs = append(cArgs, "--url", databaseURL)
			}

			if format != "" {
				cArgs = append(cArgs, "--format", format)
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
	cmd.Flags().StringVarP(&configFile, "config", "c", "", "atlas config file (default: atlas.hcl)")
	cmd.Flags().StringVarP(&env, "env", "e", "", "atlas env to use from config file")
	cmd.Flags().StringVar(&format, "format", "", "output format (e.g. sql)")

	return cmd
}
