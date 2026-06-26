package commands

import (
	"github.com/ncobase/cli/commands/create"
	initcmd "github.com/ncobase/cli/commands/init"
	"github.com/ncobase/cli/commands/migrate"
	"github.com/ncobase/cli/commands/schema"

	"github.com/spf13/cobra"
)

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
