package commands

import (
	"fmt"
	"os"

	"github.com/ncobase/cli/commands/create"
	"github.com/ncobase/cli/commands/migrate"
	"github.com/ncobase/cli/version"

	"github.com/spf13/cobra"
)

// NewStartCommand creates the start command
func NewStartCommand() *cobra.Command {
	var configFile string

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the NCore server",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("This CLI generates code projects.")
			fmt.Println("To run servers, use generated project code.")
			fmt.Println("See github.com/ncobase/ncore for library functionality.")
			return nil
		},
	}

	cmd.Flags().StringVarP(&configFile, "config", "c", "config.yaml", "config file path")
	return cmd
}

// NewPluginCommand creates the plugin management command
func NewPluginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Plugin management commands",
	}

	cmd.AddCommand(
		newPluginListCommand(),
		newPluginInstallCommand(),
	)

	return cmd
}

func newPluginListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List installed plugins",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("This CLI generates code projects.")
			fmt.Println("Plugin management belongs in generated projects.")
			fmt.Println("Use 'nco create' to generate projects.")
			return nil
		},
	}
}

func newPluginInstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install [name]",
		Short: "Install a plugin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("This CLI generates code projects.")
			fmt.Println("Plugin installation belongs in generated projects.")
			fmt.Println("Use 'nco create' to generate projects.")
			return nil
		},
	}
	return cmd
}

// NewDocsCommand creates the documentation command
func NewDocsCommand() *cobra.Command {
	var format string
	var output string

	cmd := &cobra.Command{
		Use:   "docs",
		Short: "Generate documentation",
		RunE: func(cmd *cobra.Command, args []string) error {
			var content string
			switch format {
			case "markdown":
				content = "# NCore API Documentation\n\n"
			case "json":
				content = `{"swagger": "2.0", "info": {"title": "NCore API", "version": "1.0"}}`
			default:
				return fmt.Errorf("unsupported format: %s", format)
			}

			if output == "" {
				fmt.Println(content)
				return nil
			}

			return os.WriteFile(output, []byte(content), 0644)
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "markdown", "documentation format (markdown or json)")
	cmd.Flags().StringVarP(&output, "output", "o", "", "output file path")
	return cmd
}

// NewVersionCommand creates the version command
func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			info := version.GetVersionInfo()
			fmt.Println("Version:", info.Version)
			fmt.Println("Built At:", info.BuiltAt)
		},
	}
}

// NewCreateCommand creates the extension generation command
func NewCreateCommand() *cobra.Command {
	return create.NewCommand()
}

// NewMigrateCommand creates the migrate command
func NewMigrateCommand() *cobra.Command {
	return migrate.NewCommand()
}
