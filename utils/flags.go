package utils

import "github.com/spf13/cobra"

// AddCommonFlags adds common flags shared between init and create commands
func AddCommonFlags(cmd *cobra.Command) {
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
}
