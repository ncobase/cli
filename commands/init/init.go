package initcmd

import (
	"github.com/ncobase/cli/generator"
	"github.com/spf13/cobra"
)

// NewCommand creates a new init command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [name]",
		Short: "Initialize a new project",
		Long:  `Initialize a new project with standard structure.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			opts := generator.DefaultOptions()
			opts.Name = name
			opts.Type = "direct"
			opts.Standalone = true

			// Get flags
			opts.ModuleName, _ = cmd.Flags().GetString("module")
			opts.UseMongo, _ = cmd.Flags().GetBool("use-mongo")
			opts.UseEnt, _ = cmd.Flags().GetBool("use-ent")
			opts.UseGorm, _ = cmd.Flags().GetBool("use-gorm")
			opts.WithTest, _ = cmd.Flags().GetBool("with-test")

			// Init command always generates cmd directory
			opts.WithCmd = true

			return generator.Generate(opts)
		},
	}

	// add flags
	cmd.Flags().StringP("module", "m", "", "Go module name")
	cmd.Flags().Bool("use-mongo", false, "use MongoDB")
	cmd.Flags().Bool("use-ent", false, "use Ent as ORM")
	cmd.Flags().Bool("use-gorm", false, "use Gorm as ORM")
	cmd.Flags().Bool("with-test", false, "generate test files")

	return cmd
}
