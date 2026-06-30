package initcmd

import (
	"fmt"

	"github.com/ncobase/cli/commands/internal/generation"
	"github.com/ncobase/cli/generator"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [name]",
		Short: "Initialize a new ncore application",
		Long: `Initialize a new ncore application with complete project structure.

Examples:
  nco init myapp
  nco init myapp --type modular --db postgres --use-redis --use-meilisearch
  nco init myapp --path ./apps --use-ent --db postgres
  nco init myapp --use-redis --use-kafka`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return fmt.Errorf("init requires one application name")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			opts := generator.DefaultOptions()
			opts.Name = args[0]
			opts.Type = "direct"
			opts.Standalone = true
			opts.WithCmd = true

			output, err := generation.ReadFlags(cmd, opts, generation.FlagConfig{Path: true, Type: true})
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

	generation.AddFlags(cmd, generation.FlagConfig{Path: true, Type: true})

	return cmd
}
