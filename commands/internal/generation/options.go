package generation

import (
	"encoding/json"
	"fmt"

	"github.com/ncobase/cli/generator"
	"github.com/spf13/cobra"
)

const (
	OutputText = "text"
	OutputJSON = "json"
)

// FlagConfig controls which generation flags are attached to a command.
type FlagConfig struct {
	Path    bool
	Type    bool
	WithCmd bool
	Group   bool
}

// AddFlags registers generation flags used by init and create commands.
func AddFlags(cmd *cobra.Command, cfg FlagConfig) {
	if cfg.Path {
		cmd.Flags().StringP("path", "p", "", "output path (defaults to current directory)")
	}
	if cfg.Type {
		cmd.Flags().StringP("type", "t", "service", "init type: service or modular")
	}
	cmd.Flags().StringP("module", "m", "", "Go module name")

	cmd.Flags().Bool("use-mongo", false, "use MongoDB")
	cmd.Flags().Bool("use-ent", false, "use Ent as ORM for SQL databases")
	cmd.Flags().Bool("use-gorm", false, "use GORM as ORM for SQL databases")
	cmd.Flags().String("db", "", "database driver: postgres, mysql, sqlite, mongodb")

	cmd.Flags().Bool("use-redis", false, "include Redis driver for caching and queuing")
	cmd.Flags().Bool("use-elastic", false, "include Elasticsearch driver for search")
	cmd.Flags().Bool("use-opensearch", false, "include OpenSearch driver for search")
	cmd.Flags().Bool("use-meilisearch", false, "include Meilisearch driver for search")

	cmd.Flags().Bool("use-kafka", false, "include Kafka driver for messaging")
	cmd.Flags().Bool("use-rabbitmq", false, "include RabbitMQ driver for messaging")

	cmd.Flags().Bool("use-s3", false, "include AWS S3 storage driver")
	cmd.Flags().Bool("use-minio", false, "include MinIO storage driver")
	cmd.Flags().Bool("use-aliyun", false, "include Aliyun OSS storage driver")

	cmd.Flags().Bool("with-grpc", false, "generate gRPC service support")
	cmd.Flags().Bool("with-tracing", false, "generate OpenTelemetry tracing support")
	cmd.Flags().Bool("with-test", false, "generate test files")
	if cfg.WithCmd {
		cmd.Flags().Bool("with-cmd", false, "generate cmd directory and runnable service wiring")
	}
	if cfg.Group {
		cmd.Flags().String("group", "", "optional domain group name")
	}

	cmd.Flags().Bool("dry-run", false, "print the generation plan without writing files")
	cmd.Flags().String("output", OutputText, "output format: text or json")
}

// ReadFlags copies command flag values into generator options and returns the output format.
func ReadFlags(cmd *cobra.Command, opts *generator.Options, cfg FlagConfig) (string, error) {
	var err error
	if cfg.Path {
		opts.OutputPath, err = cmd.Flags().GetString("path")
		if err != nil {
			return "", err
		}
	}
	if opts.ModuleName, err = cmd.Flags().GetString("module"); err != nil {
		return "", err
	}
	if cfg.Type {
		if opts.ProjectType, err = cmd.Flags().GetString("type"); err != nil {
			return "", err
		}
	}
	if opts.UseMongo, err = cmd.Flags().GetBool("use-mongo"); err != nil {
		return "", err
	}
	if opts.UseEnt, err = cmd.Flags().GetBool("use-ent"); err != nil {
		return "", err
	}
	if opts.UseGorm, err = cmd.Flags().GetBool("use-gorm"); err != nil {
		return "", err
	}
	if opts.DBDriver, err = cmd.Flags().GetString("db"); err != nil {
		return "", err
	}
	if opts.UseRedis, err = cmd.Flags().GetBool("use-redis"); err != nil {
		return "", err
	}
	if opts.UseElastic, err = cmd.Flags().GetBool("use-elastic"); err != nil {
		return "", err
	}
	if opts.UseOpenSearch, err = cmd.Flags().GetBool("use-opensearch"); err != nil {
		return "", err
	}
	if opts.UseMeili, err = cmd.Flags().GetBool("use-meilisearch"); err != nil {
		return "", err
	}
	if opts.UseKafka, err = cmd.Flags().GetBool("use-kafka"); err != nil {
		return "", err
	}
	if opts.UseRabbitMQ, err = cmd.Flags().GetBool("use-rabbitmq"); err != nil {
		return "", err
	}
	if opts.UseS3Storage, err = cmd.Flags().GetBool("use-s3"); err != nil {
		return "", err
	}
	if opts.UseMinio, err = cmd.Flags().GetBool("use-minio"); err != nil {
		return "", err
	}
	if opts.UseAliyun, err = cmd.Flags().GetBool("use-aliyun"); err != nil {
		return "", err
	}
	if opts.WithGRPC, err = cmd.Flags().GetBool("with-grpc"); err != nil {
		return "", err
	}
	if opts.WithTracing, err = cmd.Flags().GetBool("with-tracing"); err != nil {
		return "", err
	}
	if opts.WithTest, err = cmd.Flags().GetBool("with-test"); err != nil {
		return "", err
	}
	if cfg.WithCmd {
		if opts.WithCmd, err = cmd.Flags().GetBool("with-cmd"); err != nil {
			return "", err
		}
	}
	if cfg.Group {
		if opts.Group, err = cmd.Flags().GetString("group"); err != nil {
			return "", err
		}
	}
	if opts.DryRun, err = cmd.Flags().GetBool("dry-run"); err != nil {
		return "", err
	}

	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return "", err
	}
	if output != OutputText && output != OutputJSON {
		return "", fmt.Errorf("unsupported output format %q; supported formats are text and json", output)
	}

	return output, nil
}

// WriteResult writes a generation result in the requested format.
func WriteResult(cmd *cobra.Command, result *generator.Result, output string) error {
	switch output {
	case OutputJSON:
		encoder := json.NewEncoder(cmd.OutOrStdout())
		encoder.SetIndent("", "  ")
		return encoder.Encode(result)
	case OutputText:
		_, err := fmt.Fprintln(cmd.OutOrStdout(), result.Text())
		return err
	default:
		return fmt.Errorf("unsupported output format %q", output)
	}
}
