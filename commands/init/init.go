package initcmd

import (
	"github.com/ncobase/cli/generator"
	"github.com/ncobase/cli/utils"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [name]",
		Short: "Initialize a new standalone application",
		Long: `Initialize a new standalone application with complete project structure.

Examples:
  nco init myapp
  nco init myapp --use-ent --db postgres
  nco init myapp --use-redis --use-kafka`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := generator.DefaultOptions()
			opts.Name = args[0]
			opts.Type = "direct"
			opts.Standalone = true
			opts.WithCmd = true

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
			opts.WithGRPC, _ = cmd.Flags().GetBool("with-grpc")
			opts.WithTracing, _ = cmd.Flags().GetBool("with-tracing")

			return generator.Generate(opts)
		},
	}

	cmd.Flags().StringP("module", "m", "", "Go module name (e.g., github.com/username/project)")
	utils.AddCommonFlags(cmd)

	return cmd
}
