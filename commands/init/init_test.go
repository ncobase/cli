package initcmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewCommand_MetadataAndFlags(t *testing.T) {
	cmd := NewCommand()

	if cmd.Use != "init [name]" {
		t.Fatalf("unexpected use field: %q", cmd.Use)
	}

	requiredFlags := []string{
		"path", "module", "use-mongo", "use-ent", "use-gorm", "db",
		"use-redis", "use-elastic", "use-opensearch", "use-meilisearch",
		"use-kafka", "use-rabbitmq", "use-s3", "use-minio", "use-aliyun",
		"type", "with-grpc", "with-tracing", "with-test", "dry-run", "output",
	}

	for _, name := range requiredFlags {
		if cmd.Flags().Lookup(name) == nil {
			t.Fatalf("expected init flag %q to exist", name)
		}
	}
}

func TestNewCommand_NoArgsShowsHelp(t *testing.T) {
	cmd := NewCommand()
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected init help flow to succeed, got error: %v", err)
	}
	if !strings.Contains(output.String(), "Usage:") || !strings.Contains(output.String(), "nco init myapp") {
		t.Fatalf("expected init help output, got %q", output.String())
	}
	if strings.Contains(output.String(), "accepts") {
		t.Fatalf("expected init help instead of raw arg error, got %q", output.String())
	}
}

func TestNewCommand_TooManyArgsReturnsFriendlyError(t *testing.T) {
	cmd := NewCommand()
	cmd.SetArgs([]string{"one", "two"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected init command to reject extra arguments")
	}
	if err.Error() != "init requires one application name" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewCommand_ModularDryRunJSON(t *testing.T) {
	cmd := NewCommand()
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetArgs([]string{"product", "--type", "modular", "--dry-run", "--output", "json"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected modular dry-run to succeed, got error: %v", err)
	}

	for _, expected := range []string{
		`"project_type": "modular"`,
		`"core/doc.go"`,
		`"biz/doc.go"`,
		`"internal/server/http.go"`,
	} {
		if !strings.Contains(output.String(), expected) {
			t.Fatalf("expected output to contain %s, got %q", expected, output.String())
		}
	}
}
