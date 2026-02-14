package initcmd

import (
	"github.com/ncobase/cli/generator"
	"github.com/spf13/cobra"
)

// NewCommand creates a new init command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [name]",
		Short: "Initialize a new standalone application",
		Long: `Initialize a new standalone application with complete project structure.

This command creates a ready-to-run Go application with:
  - Complete directory structure (cmd, data, handler, service, etc.)
  - Configuration files (config.yaml, .gitignore, go.mod)
  - Optional ORM support (Ent, GORM, or MongoDB)
  - Optional data source drivers (Redis, Elasticsearch, Kafka, etc.)
  - Optional storage drivers (S3, MinIO, Aliyun OSS)

Examples:
  nco init myapp                           # Basic application
  nco init myapp --use-ent --db postgres   # With Ent ORM and PostgreSQL
  nco init myapp --use-redis --use-kafka   # With Redis and Kafka
  nco init myapp -m github.com/me/myapp    # Custom module name

For creating extensions within an existing project, use 'nco create' instead.`,
		Args: cobra.ExactArgs(1),
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
			opts.WithTest, _ = cmd.Flags().GetBool("with-test")

			// Init command always generates cmd directory
			opts.WithCmd = true

			return generator.Generate(opts)
		},
	}

	// add flags
	cmd.Flags().StringP("module", "m", "", "Go module name (e.g., github.com/username/project)")

	// ORM options (mutually exclusive)
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

	// Other options
	cmd.Flags().Bool("with-test", false, "generate test files (unit, integration, e2e)")

	return cmd
}
