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
	cmd.Flags().StringP("module", "m", "", "Go module name")
	cmd.Flags().Bool("use-mongo", false, "use MongoDB")
	cmd.Flags().Bool("use-ent", false, "use Ent as ORM")
	cmd.Flags().Bool("use-gorm", false, "use Gorm as ORM")
	cmd.Flags().String("db", "", "database driver (postgres, mysql, sqlite, mongodb, neo4j)")
	cmd.Flags().Bool("use-redis", false, "include Redis driver")
	cmd.Flags().Bool("use-elastic", false, "include Elasticsearch driver")
	cmd.Flags().Bool("use-opensearch", false, "include OpenSearch driver")
	cmd.Flags().Bool("use-meilisearch", false, "include Meilisearch driver")
	cmd.Flags().Bool("use-kafka", false, "include Kafka driver")
	cmd.Flags().Bool("use-rabbitmq", false, "include RabbitMQ driver")
	cmd.Flags().Bool("use-s3", false, "include S3 storage driver")
	cmd.Flags().Bool("use-minio", false, "include MinIO storage driver")
	cmd.Flags().Bool("use-aliyun", false, "include Aliyun OSS storage driver")
	cmd.Flags().Bool("with-test", false, "generate test files")

	return cmd
}
