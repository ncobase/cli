package commands

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	originalStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}

	os.Stdout = writer
	fn()

	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close pipe writer: %v", err)
	}
	os.Stdout = originalStdout

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, reader); err != nil {
		t.Fatalf("failed to read stdout: %v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("failed to close pipe reader: %v", err)
	}

	return buf.String()
}

func TestNewRootCmd_Subcommands(t *testing.T) {
	cmd := NewRootCmd()

	if cmd.Use != filepath.Base(os.Args[0]) {
		t.Fatalf("unexpected root use value: %q", cmd.Use)
	}

	expected := map[string]bool{
		"version": false,
		"init":    false,
		"create":  false,
		"migrate": false,
		"schema":  false,
	}

	for _, sub := range cmd.Commands() {
		if _, ok := expected[sub.Name()]; ok {
			expected[sub.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected subcommand %q to be registered", name)
		}
	}
}

func TestNewVersionCommand_Run(t *testing.T) {
	cmd := NewVersionCommand()
	output := captureStdout(t, func() {
		cmd.Run(cmd, nil)
	})

	if !strings.Contains(output, "Version:") {
		t.Fatalf("expected output to contain version line, got %q", output)
	}
	if !strings.Contains(output, "Built At:") {
		t.Fatalf("expected output to contain build time line, got %q", output)
	}
}
