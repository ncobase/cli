package commands

import (
	"fmt"

	"github.com/ncobase/cli/commands/create"
	initcmd "github.com/ncobase/cli/commands/init"
	"github.com/ncobase/cli/commands/migrate"
	"github.com/ncobase/cli/commands/schema"
	"github.com/ncobase/cli/version"

	"github.com/spf13/cobra"
)

// NewVersionCommand creates the version command
func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			info := version.GetVersionInfo()
			out := cmd.OutOrStdout()
			fmt.Fprintln(out, "Version:", info.Version)
			fmt.Fprintln(out, "Built At:", info.BuiltAt)
			return nil
		},
	}
}

// NewCreateCommand creates the extension generation command
func NewCreateCommand() *cobra.Command {
	return create.NewCommand()
}

// NewInitCommand creates the init command
func NewInitCommand() *cobra.Command {
	return initcmd.NewCommand()
}

// NewMigrateCommand creates the migrate command
func NewMigrateCommand() *cobra.Command {
	return migrate.NewCommand()
}

// NewSchemaCommand creates the schema command
func NewSchemaCommand() *cobra.Command {
	return schema.NewCommand()
}
