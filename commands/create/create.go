package create

import (
	"strings"

	"github.com/ncobase/cli/generator"

	"github.com/spf13/cobra"
)

// knownTypes is a map of known extension types
var knownTypes = map[string]string{
	"core":     "core",
	"business": "business",
	"plugin":   "plugin",
}

// NewCommand creates a new create command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create [extension type or custom directory] name",
		Aliases: []string{"gen", "generate"},
		Short:   "Create new extension components",
		Long: `Create new extensions within an existing ncobase project.

Extension types:
  core      - Core domain extensions (fundamental business logic)
  business  - Business domain extensions (application-specific logic)
  plugin    - Plugin extensions (optional features)
  [custom]  - Custom directory name

Examples:
  nco create core auth          # Create core/auth extension
  nco create business crm       # Create business/crm extension
  nco create plugin payment     # Create plugin/payment extension
  nco create myext user         # Create myext/user in custom directory

For creating standalone applications, use 'nco init' instead.`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := generator.DefaultOptions()

			var dir, name string
			if len(args) == 1 {
				// If only one argument, check if it's a known type
				firstArg := strings.ToLower(args[0])
				if _, ok := knownTypes[firstArg]; ok {
					// Is a known type, show help
					cmd.Help()
					return nil
				}

				// Not a known type, assume it's the name and create directly in current directory
				name = args[0]

				// Set options
				opts.Name = name
				opts.Type = "direct" // New type for direct creation

				// Get flags
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
			}

			// If two arguments, use the first as the directory and the second as the name
			dir = args[0]
			name = args[1]

			// Check if the directory is a known type
			if extType, ok := knownTypes[strings.ToLower(dir)]; ok {
				// Is a known type
				switch extType {
				case "core":
					return newCoreCommand().RunE(cmd, []string{name})
				case "business":
					return newBusinessCommand().RunE(cmd, []string{name})
				case "plugin":
					return newPluginCommand().RunE(cmd, []string{name})
				}
			}

			// Not a known type, assume it's a custom directory
			opts.Name = name
			opts.Type = "custom"
			opts.CustomDir = dir

			// Get flags
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

	// add subcommands
	cmd.AddCommand(
		newCoreCommand(),
		newBusinessCommand(),
		newPluginCommand(),
	)

	// add flags
	cmd.Flags().StringP("path", "p", "", "output path (defaults to current directory)")
	cmd.Flags().StringP("module", "m", "", "Go module name (defaults to current module)")

	// ORM options
	cmd.Flags().Bool("use-mongo", false, "use MongoDB as database")
	cmd.Flags().Bool("use-ent", false, "use Ent as ORM (for SQL databases)")
	cmd.Flags().Bool("use-gorm", false, "use GORM as ORM (for SQL databases)")
	cmd.Flags().String("db", "", "database driver: postgres, mysql, sqlite, mongodb, neo4j")

	// Data source drivers
	cmd.Flags().Bool("use-redis", false, "include Redis driver for caching/queuing")
	cmd.Flags().Bool("use-elastic", false, "include Elasticsearch driver for search")
	cmd.Flags().Bool("use-opensearch", false, "include OpenSearch driver for search")
	cmd.Flags().Bool("use-meilisearch", false, "include Meilisearch driver for search")

	// Message queue drivers
	cmd.Flags().Bool("use-kafka", false, "include Kafka driver for messaging")
	cmd.Flags().Bool("use-rabbitmq", false, "include RabbitMQ driver for messaging")

	// Storage drivers
	cmd.Flags().Bool("use-s3", false, "include AWS S3 storage driver")
	cmd.Flags().Bool("use-minio", false, "include MinIO storage driver")
	cmd.Flags().Bool("use-aliyun", false, "include Aliyun OSS storage driver")

	// Service options
	cmd.Flags().Bool("with-grpc", false, "generate gRPC service support")
	cmd.Flags().Bool("with-tracing", false, "generate OpenTelemetry tracing support")

	// Other options
	cmd.Flags().Bool("with-test", false, "generate test files")
	cmd.Flags().Bool("with-cmd", false, "generate cmd directory with main.go for testing")
	cmd.Flags().String("group", "", "optional domain group name")

	return cmd
}
