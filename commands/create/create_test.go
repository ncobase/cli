package create

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func captureOutput(t *testing.T, fn func() error) error {
	t.Helper()

	originalStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}

	os.Stdout = writer
	runErr := fn()

	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close pipe writer: %v", err)
	}
	os.Stdout = originalStdout

	if _, err := io.Copy(io.Discard, reader); err != nil {
		t.Fatalf("failed to drain stdout: %v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("failed to close pipe reader: %v", err)
	}

	return runErr
}

func TestNewCommand_HasExpectedSubcommandsAndFlags(t *testing.T) {
	cmd := NewCommand()

	expectedSubcommands := map[string]bool{
		"core":     false,
		"business": false,
		"plugin":   false,
	}

	for _, sub := range cmd.Commands() {
		if _, ok := expectedSubcommands[sub.Name()]; ok {
			expectedSubcommands[sub.Name()] = true
		}
	}

	for name, found := range expectedSubcommands {
		if !found {
			t.Fatalf("expected subcommand %q to be registered", name)
		}
	}

	requiredFlags := []string{
		"path", "module", "use-mongo", "use-ent", "use-gorm", "db",
		"use-redis", "use-elastic", "use-opensearch", "use-meilisearch",
		"use-kafka", "use-rabbitmq", "use-s3", "use-minio", "use-aliyun",
		"with-grpc", "with-tracing", "with-test", "with-cmd", "group",
		"dry-run", "output",
	}

	for _, flagName := range requiredFlags {
		if cmd.Flags().Lookup(flagName) == nil {
			t.Fatalf("expected flag %q to exist", flagName)
		}
	}
}

func TestCreateSubcommandsHaveGenerationFlags(t *testing.T) {
	cmd := NewCommand()
	for _, name := range []string{"core", "business", "plugin"} {
		sub, _, err := cmd.Find([]string{name})
		if err != nil {
			t.Fatalf("failed to find %s subcommand: %v", name, err)
		}
		for _, flagName := range []string{"path", "module", "with-cmd", "group", "dry-run", "output"} {
			if sub.Flags().Lookup(flagName) == nil {
				t.Fatalf("expected %s subcommand to include flag %q", name, flagName)
			}
		}
	}
}

func TestCreateKnownTypeWithSingleArgumentShowsHelp(t *testing.T) {
	cmd := NewCommand()
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetErr(&output)

	if err := cmd.RunE(cmd, []string{"core"}); err != nil {
		t.Fatalf("expected help flow to succeed, got error: %v", err)
	}
	if output.Len() == 0 {
		t.Fatal("expected help output for known type input")
	}
}

func TestCreateDirectGeneration(t *testing.T) {
	tmpDir := t.TempDir()
	cmd := NewCommand()
	cmd.SetArgs([]string{"sampleext", "--path", tmpDir})

	if err := captureOutput(t, cmd.Execute); err != nil {
		t.Fatalf("direct create command failed: %v", err)
	}

	base := filepath.Join(tmpDir, "sampleext")
	expectedFiles := []string{
		filepath.Join(base, "sampleext.go"),
		filepath.Join(base, "data/data.go"),
		filepath.Join(base, "handler/provider.go"),
		filepath.Join(base, "service/provider.go"),
	}

	for _, file := range expectedFiles {
		if _, err := os.Stat(file); err != nil {
			t.Fatalf("expected generated file %q, got error: %v", file, err)
		}
	}
}

func TestCreateCustomDirectoryGeneration(t *testing.T) {
	tmpDir := t.TempDir()
	cmd := NewCommand()
	cmd.SetArgs([]string{"customdir", "sampleext", "--path", tmpDir, "--with-test"})

	if err := captureOutput(t, cmd.Execute); err != nil {
		t.Fatalf("custom create command failed: %v", err)
	}

	base := filepath.Join(tmpDir, "customdir", "sampleext")
	expectedFiles := []string{
		filepath.Join(base, "sampleext.go"),
		filepath.Join(base, "tests/ext_test.go"),
		filepath.Join(base, "tests/handler_test.go"),
		filepath.Join(base, "tests/service_test.go"),
	}

	for _, file := range expectedFiles {
		if _, err := os.Stat(file); err != nil {
			t.Fatalf("expected generated file %q, got error: %v", file, err)
		}
	}
}

func TestCreateDryRunJSONDoesNotWriteFiles(t *testing.T) {
	tmpDir := t.TempDir()
	cmd := NewCommand()
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetArgs([]string{"core", "audit", "--path", tmpDir, "--dry-run", "--output", "json", "--use-gorm", "--db", "sqlite", "--with-cmd"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("dry-run create command failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "core", "audit")); !os.IsNotExist(err) {
		t.Fatalf("expected dry-run not to create files, stat error: %v", err)
	}
	for _, expected := range []string{`"dry_run": true`, `"base_path"`, `"module_requirements"`} {
		if !bytes.Contains(output.Bytes(), []byte(expected)) {
			t.Fatalf("expected json output to contain %s, got %s", expected, output.String())
		}
	}
}
