package create

import (
	"strings"

	"github.com/ncobase/cli/generator"
	"github.com/ncobase/cli/utils"
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
			opts := generator.DefaultOptions()

			var dir, name string
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
					switch extType {
					case "core":
						return newCoreCommand().RunE(cmd, []string{name})
					case "business":
						return newBusinessCommand().RunE(cmd, []string{name})
					case "plugin":
						return newPluginCommand().RunE(cmd, []string{name})
					}
				}
				opts.Type = "custom"
				opts.CustomDir = dir
			}

			opts.Name = name
			opts.ModuleName, _ = cmd.Flags().GetString("module")
			opts.OutputPath, _ = cmd.Flags().GetString("path")
			opts.UseMongo, _ = cmd.Flags().GetBool("use-mongo")
			opts.UseEnt, _ = cmd.Flags().GetBool("use-ent")
			opts.UseGorm, _ = cmd.Flags().GetBool("use-gorm")
			opts.DBDriver, _ = cmd.Flags().GetString("db")
			opts.UseRedis, _ = cmd.Flags().GetBool("use-redis")
			opts.UseElastic, _ = cmd.Flags().GetBool("use-elastic")
			opts.UseOpenSearch, _ = cmd.Flags().GetBool("use-opensearch")
			opts.UseMeili, _ = cmd.Flags().GetBool("use-meilisearch")
			opts.UseKafka, _ = cmd.Flags().GetBool("use-kafka")
			opts.UseRabbitMQ, _ = cmd.Flags().GetBool("use-rabbitmq")
			opts.UseS3Storage, _ = cmd.Flags().GetBool("use-s3")
			opts.UseMinio, _ = cmd.Flags().GetBool("use-minio")
			opts.UseAliyun, _ = cmd.Flags().GetBool("use-aliyun")
			opts.WithCmd, _ = cmd.Flags().GetBool("with-cmd")
			opts.WithTest, _ = cmd.Flags().GetBool("with-test")
			opts.WithGRPC, _ = cmd.Flags().GetBool("with-grpc")
			opts.WithTracing, _ = cmd.Flags().GetBool("with-tracing")
			opts.Group, _ = cmd.Flags().GetString("group")

			return generator.Generate(opts)
		},
	}

	cmd.AddCommand(newCoreCommand(), newBusinessCommand(), newPluginCommand())

	cmd.Flags().StringP("path", "p", "", "output path (defaults to current directory)")
	cmd.Flags().StringP("module", "m", "", "Go module name (defaults to current module)")
	cmd.Flags().Bool("with-cmd", false, "generate cmd directory with main.go for testing")
	cmd.Flags().String("group", "", "optional domain group name")
	utils.AddCommonFlags(cmd)

	return cmd
}
