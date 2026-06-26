package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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
	var output bytes.Buffer
	cmd.SetOut(&output)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("version command failed: %v", err)
	}

	if !strings.Contains(output.String(), "Version:") {
		t.Fatalf("expected output to contain version line, got %q", output.String())
	}
	if !strings.Contains(output.String(), "Built At:") {
		t.Fatalf("expected output to contain build time line, got %q", output.String())
	}
}

func TestNewVersionCommand_Verbose(t *testing.T) {
	cmd := NewVersionCommand()
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetArgs([]string{"--verbose"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("version command failed: %v", err)
	}

	for _, expected := range []string{"Version:", "Branch:", "Revision:", "Go Version:"} {
		if !strings.Contains(output.String(), expected) {
			t.Fatalf("expected verbose output to contain %q, got %q", expected, output.String())
		}
	}
}

func TestNewVersionCommand_JSON(t *testing.T) {
	cmd := NewVersionCommand()
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetArgs([]string{"--json"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("version command failed: %v", err)
	}

	for _, expected := range []string{`"version"`, `"built_at"`, `"go_version"`} {
		if !strings.Contains(output.String(), expected) {
			t.Fatalf("expected json output to contain %q, got %q", expected, output.String())
		}
	}
}
