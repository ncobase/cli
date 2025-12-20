package migrate

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func newHashCommand() *cobra.Command {
	var dirPath string

	cmd := &cobra.Command{
		Use:   "hash",
		Short: "Re-hash the migration directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Default values
			if dirPath == "" {
				dirPath = "file://data/ent/migrate/migrations"
			}

			args = []string{"migrate", "hash", "--dir", dirPath}

			fmt.Printf("Running: atlas %v\n", args)

			c := exec.Command("atlas", args...)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin

			return c.Run()
		},
	}

	cmd.Flags().StringVar(&dirPath, "dir", "", "migration directory (default: file://data/ent/migrate/migrations)")

	return cmd
}
