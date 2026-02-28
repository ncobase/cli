package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestWriteTemplateFile(t *testing.T) {
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "nested", "hello.txt")

	err := WriteTemplateFile(target, "hello {{.Name}}", map[string]string{"Name": "ncobase"})
	if err != nil {
		t.Fatalf("WriteTemplateFile() error = %v", err)
	}

	content, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(content) != "hello ncobase" {
		t.Fatalf("unexpected file content: %q", string(content))
	}
}

func TestWriteTemplateFile_InvalidTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "broken.txt")

	err := WriteTemplateFile(target, "{{ if .Name }}", map[string]string{"Name": "ncobase"})
	if err == nil {
		t.Fatal("expected template parse error, got nil")
	}
}

func TestValidateName(t *testing.T) {
	valid := []string{"sample", "sample_1", "SampleName", "sample-name"}
	for _, name := range valid {
		if !ValidateName(name) {
			t.Fatalf("expected name %q to be valid", name)
		}
	}

	invalid := []string{"", "1sample", "sample name", "sample.name", "@sample"}
	for _, name := range invalid {
		if ValidateName(name) {
			t.Fatalf("expected name %q to be invalid", name)
		}
	}
}

func TestPathAndFileExists(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "a.txt")

	if PathExists(file) || FileExists(file) {
		t.Fatal("expected non-existing file to return false")
	}

	if err := os.WriteFile(file, []byte("ok"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if !PathExists(file) || !FileExists(file) {
		t.Fatal("expected existing file to return true")
	}
	if !PathExists(tmpDir) {
		t.Fatal("expected existing directory to return true")
	}
}

func TestGetPlatformExt(t *testing.T) {
	got := GetPlatformExt()
	if runtime.GOOS == "windows" && got != ".exe" {
		t.Fatalf("expected .exe on windows, got %q", got)
	}
	if runtime.GOOS != "windows" && got != "" {
		t.Fatalf("expected empty extension on non-windows, got %q", got)
	}
}
