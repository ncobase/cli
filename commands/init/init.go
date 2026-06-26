package initcmd

import (
	"github.com/ncobase/cli/commands/internal/generation"
	"github.com/ncobase/cli/generator"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [name]",
		Short: "Initialize a new standalone application",
		Long: `Initialize a new standalone application with complete project structure.

		Examples:
  nco init myapp
  nco init myapp --path ./apps --use-ent --db postgres
  nco init myapp --use-redis --use-kafka`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := generator.DefaultOptions()
			opts.Name = args[0]
			opts.Type = "direct"
			opts.Standalone = true
			opts.WithCmd = true

			output, err := generation.ReadFlags(cmd, opts, generation.FlagConfig{Path: true})
			if err != nil {
				return err
			}

			result, err := generator.Generate(opts)
			if err != nil {
				return err
			}
			return generation.WriteResult(cmd, result, output)
		},
	}

	generation.AddFlags(cmd, generation.FlagConfig{Path: true})

	return cmd
}
