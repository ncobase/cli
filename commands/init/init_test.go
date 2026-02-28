package initcmd

import "testing"

func TestNewCommand_MetadataAndFlags(t *testing.T) {
	cmd := NewCommand()

	if cmd.Use != "init [name]" {
		t.Fatalf("unexpected use field: %q", cmd.Use)
	}

	requiredFlags := []string{
		"module", "use-mongo", "use-ent", "use-gorm", "db",
		"use-redis", "use-elastic", "use-opensearch", "use-meilisearch",
		"use-kafka", "use-rabbitmq", "use-s3", "use-minio", "use-aliyun",
		"with-grpc", "with-tracing", "with-test",
	}

	for _, name := range requiredFlags {
		if cmd.Flags().Lookup(name) == nil {
			t.Fatalf("expected init flag %q to exist", name)
		}
	}
}
