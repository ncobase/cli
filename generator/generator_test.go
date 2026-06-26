package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateDryRunBuildsPlanWithoutWritingFiles(t *testing.T) {
	tmpDir := t.TempDir()
	opts := DefaultOptions()
	opts.Name = "sample"
	opts.Type = "direct"
	opts.Standalone = true
	opts.WithCmd = true
	opts.OutputPath = tmpDir
	opts.UseGorm = true
	opts.DBDriver = "sqlite"
	opts.DryRun = true

	result, err := Generate(opts)
	if err != nil {
		t.Fatalf("dry-run generation failed: %v", err)
	}
	if !result.DryRun || result.Applied {
		t.Fatalf("unexpected dry-run result: %#v", result)
	}
	if result.Plan.BasePath != filepath.Join(tmpDir, "sample") {
		t.Fatalf("unexpected base path: %q", result.Plan.BasePath)
	}
	if len(result.Plan.Files) == 0 || len(result.Plan.Directories) == 0 {
		t.Fatalf("expected dry-run plan to include files and directories: %#v", result.Plan)
	}
	if _, err := os.Stat(result.Plan.BasePath); !os.IsNotExist(err) {
		t.Fatalf("expected dry-run not to create base path, stat error: %v", err)
	}
}

func TestGenerateDryRunReportsConflicts(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.Mkdir(filepath.Join(tmpDir, "sample"), 0755); err != nil {
		t.Fatalf("failed to create conflict directory: %v", err)
	}

	opts := DefaultOptions()
	opts.Name = "sample"
	opts.Type = "direct"
	opts.Standalone = true
	opts.OutputPath = tmpDir
	opts.DryRun = true

	result, err := Generate(opts)
	if err != nil {
		t.Fatalf("dry-run conflict planning failed: %v", err)
	}
	if len(result.Plan.Conflicts) != 1 {
		t.Fatalf("expected one conflict, got %#v", result.Plan.Conflicts)
	}
}

func TestGenerateRejectsUnsupportedDatabaseDriver(t *testing.T) {
	opts := DefaultOptions()
	opts.Name = "graph"
	opts.Type = "direct"
	opts.Standalone = true
	opts.DBDriver = "neo4j"
	opts.DryRun = true

	if _, err := Generate(opts); err == nil {
		t.Fatal("expected unsupported database driver error")
	}
}

func TestGenerateRejectsMultipleStorageDrivers(t *testing.T) {
	opts := DefaultOptions()
	opts.Name = "storage"
	opts.Type = "direct"
	opts.Standalone = true
	opts.UseMinio = true
	opts.UseS3Storage = true
	opts.DryRun = true

	if _, err := Generate(opts); err == nil {
		t.Fatal("expected storage driver conflict error")
	}
}
