package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ncobase/cli/version"
	"github.com/spf13/cobra"
)

const rootHelpTemplate = `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}Usage:
  {{.UseLine}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

// NewRootCmd creates the root command
func NewRootCmd() *cobra.Command {
	var showVersion bool

	rootCmd := &cobra.Command{
		Use:           filepath.Base(os.Args[0]),
		Short:         "Scaffold Go applications with the ncore framework",
		Long:          "Ncobase CLI scaffolds Go applications and extension modules for the ncore framework.",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if showVersion {
				info := version.GetVersionInfo()
				fmt.Fprintln(cmd.OutOrStdout(), "Version:", info.Version)
				fmt.Fprintln(cmd.OutOrStdout(), "Built At:", info.BuiltAt)
				return nil
			}
			return cmd.Help()
		},
	}
	rootCmd.SetHelpTemplate(rootHelpTemplate)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "print version information")

	rootCmd.AddCommand(
		NewInitCommand(),
		NewCreateCommand(),
		NewMigrateCommand(),
		NewSchemaCommand(),
	)

	return rootCmd
}
