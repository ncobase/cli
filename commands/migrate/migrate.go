package migrate

import (
	"github.com/ncobase/cli/utils"
	"github.com/spf13/cobra"
)

// NewCommand creates a new migrate command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "migrate",
		Args:    cobra.NoArgs,
		Aliases: []string{"m"},
		Short:   "Database migration commands",
		Long:    `Manage database migrations.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return utils.CheckAndInstallAtlas()
		},
	}

	cmd.AddCommand(
		newApplyCommand(),
		newDownCommand(),
		newNewCommand(),
		newDiffCommand(),
		newHashCommand(),
		newStatusCommand(),
		newLintCommand(),
	)

	return cmd
}
