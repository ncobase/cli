package create

import (
	"github.com/ncobase/cli/generator"

	"github.com/spf13/cobra"
)

func newCoreCommand() *cobra.Command {
	opts := &generator.Options{}

	cmd := &cobra.Command{
		Use:     "core [name]",
		Aliases: []string{"c"},
		Short:   "Create a new extension in core domain",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]
			opts.Type = "core"

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
			opts.Standalone, _ = cmd.Flags().GetBool("standalone")
			opts.Group, _ = cmd.Flags().GetString("group")

			return generator.Generate(opts)
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&opts.OutputPath, "path", "p", "", "output path (defaults to current directory)")
	cmd.Flags().StringVarP(&opts.ModuleName, "module", "m", "", "Go module name (defaults to current module)")
	cmd.Flags().BoolVar(&opts.UseMongo, "use-mongo", false, "use MongoDB")
	cmd.Flags().BoolVar(&opts.UseEnt, "use-ent", false, "use Ent as ORM")
	cmd.Flags().BoolVar(&opts.UseGorm, "use-gorm", false, "use Gorm as ORM")
	cmd.Flags().StringVar(&opts.DBDriver, "db", "", "database driver (postgres, mysql, sqlite, mongodb, neo4j)")
	cmd.Flags().BoolVar(&opts.UseRedis, "use-redis", false, "include Redis driver")
	cmd.Flags().BoolVar(&opts.UseElastic, "use-elastic", false, "include Elasticsearch driver")
	cmd.Flags().BoolVar(&opts.UseOpenSearch, "use-opensearch", false, "include OpenSearch driver")
	cmd.Flags().BoolVar(&opts.UseMeili, "use-meilisearch", false, "include Meilisearch driver")
	cmd.Flags().BoolVar(&opts.UseKafka, "use-kafka", false, "include Kafka driver")
	cmd.Flags().BoolVar(&opts.UseRabbitMQ, "use-rabbitmq", false, "include RabbitMQ driver")
	cmd.Flags().BoolVar(&opts.UseS3Storage, "use-s3", false, "include S3 storage driver")
	cmd.Flags().BoolVar(&opts.UseMinio, "use-minio", false, "include MinIO storage driver")
	cmd.Flags().BoolVar(&opts.UseAliyun, "use-aliyun", false, "include Aliyun OSS storage driver")
	cmd.Flags().BoolVar(&opts.WithCmd, "with-cmd", false, "generate cmd directory with main.go")
	cmd.Flags().BoolVar(&opts.WithTest, "with-test", false, "generate test files")
	cmd.Flags().BoolVar(&opts.WithGRPC, "with-grpc", false, "generate gRPC service support")
	cmd.Flags().BoolVar(&opts.WithTracing, "with-tracing", false, "generate OpenTelemetry tracing support")
	cmd.Flags().BoolVar(&opts.Standalone, "standalone", false, "generate as standalone app without extension structure")
	cmd.Flags().StringVar(&opts.Group, "group", "", "belongs domain group (optional)")

	return cmd
}
