package schema

import (
	"github.com/ncobase/cli/utils"
	"github.com/spf13/cobra"
)

// NewCommand creates a new schema command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schema",
		Short: "Work with database schemas",
		Long:  `The schema command allows you to inspect, apply, and manage database schemas.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return utils.CheckAndInstallAtlas()
		},
	}

	cmd.AddCommand(
		newInspectCommand(),
		newApplyCommand(),
	)

	return cmd
}
