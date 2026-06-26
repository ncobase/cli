package create

import (
	"strings"

	"github.com/ncobase/cli/commands/internal/generation"
	"github.com/ncobase/cli/generator"
	"github.com/spf13/cobra"
)

var knownTypes = map[string]string{
	"core":     "core",
	"business": "business",
	"plugin":   "plugin",
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create [extension type or custom directory] name",
		Aliases: []string{"gen", "generate"},
		Short:   "Create new extension components",
		Long: `Create new extensions within an existing ncobase project.

Examples:
  nco create core auth
  nco create business crm
  nco create plugin payment`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var dir, name string
			opts := generator.DefaultOptions()
			if len(args) == 1 {
				firstArg := strings.ToLower(args[0])
				if _, ok := knownTypes[firstArg]; ok {
					cmd.Help()
					return nil
				}
				name = args[0]
				opts.Type = "direct"
			} else {
				dir = args[0]
				name = args[1]
				if extType, ok := knownTypes[strings.ToLower(dir)]; ok {
					opts.Type = extType
				} else {
					opts.Type = "custom"
					opts.CustomDir = dir
				}
			}

			opts.Name = name
			output, err := generation.ReadFlags(cmd, opts, generation.FlagConfig{Path: true, WithCmd: true, Group: true})
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

	cmd.AddCommand(newTypedCommand("core", []string{"c"}, "Create a new extension in core domain"))
	cmd.AddCommand(newTypedCommand("business", []string{"b"}, "Create a new extension in business domain"))
	cmd.AddCommand(newTypedCommand("plugin", []string{"p"}, "Create a new extension in plugin domain"))

	generation.AddFlags(cmd, generation.FlagConfig{Path: true, WithCmd: true, Group: true})

	return cmd
}

func newTypedCommand(extType string, aliases []string, short string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     extType + " [name]",
		Aliases: aliases,
		Short:   short,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := generator.DefaultOptions()
			opts.Name = args[0]
			opts.Type = extType

			output, err := generation.ReadFlags(cmd, opts, generation.FlagConfig{Path: true, WithCmd: true, Group: true})
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

	generation.AddFlags(cmd, generation.FlagConfig{Path: true, WithCmd: true, Group: true})
	return cmd
}
