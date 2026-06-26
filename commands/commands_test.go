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

	if !cmd.CompletionOptions.DisableDefaultCmd {
		t.Fatal("expected default completion command to be disabled")
	}
	if cmd.SilenceUsage != true || cmd.SilenceErrors != true {
		t.Fatal("expected root command to silence cobra usage and error output")
	}
	if cmd.Flags().Lookup("version") == nil {
		t.Fatal("expected root version flag to be registered")
	}
	if flag := cmd.Flags().ShorthandLookup("v"); flag == nil || flag.Name != "version" {
		t.Fatal("expected -v shorthand to map to version flag")
	}
}

func TestNewRootCmd_DefaultShowsHelp(t *testing.T) {
	cmd := NewRootCmd()
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("root command failed: %v", err)
	}
	if !strings.Contains(output.String(), "Usage:") || !strings.Contains(output.String(), "Available Commands:") {
		t.Fatalf("expected root command to print help, got %q", output.String())
	}
	if strings.Contains(output.String(), "completion") {
		t.Fatalf("expected help output not to include completion command, got %q", output.String())
	}
	if strings.Contains(output.String(), "\n  help") {
		t.Fatalf("expected help output not to include help command, got %q", output.String())
	}
	if strings.Contains(output.String(), "\n  version") {
		t.Fatalf("expected help output not to include version command, got %q", output.String())
	}
}

func TestNewRootCmd_VersionFlag(t *testing.T) {
	for _, arg := range []string{"--version", "-v"} {
		cmd := NewRootCmd()
		var output bytes.Buffer
		cmd.SetOut(&output)
		cmd.SetArgs([]string{arg})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("root version flag %s failed: %v", arg, err)
		}
		if !strings.Contains(output.String(), "Version:") || !strings.Contains(output.String(), "Built At:") {
			t.Fatalf("expected version output for %s, got %q", arg, output.String())
		}
	}
}

func TestNewRootCmd_UnknownCommandReturnsSingleError(t *testing.T) {
	cmd := NewRootCmd()
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetErr(&output)
	cmd.SetArgs([]string{"unknown"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected unknown command to return an error")
	}
	if output.Len() != 0 {
		t.Fatalf("expected cobra not to print duplicate error output, got %q", output.String())
	}
}

func TestNewRootCmd_VersionSubcommandIsNotRegistered(t *testing.T) {
	cmd := NewRootCmd()
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetErr(&output)
	cmd.SetArgs([]string{"version"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected version subcommand to be unavailable")
	}
	if !strings.Contains(err.Error(), `unknown command "version"`) {
		t.Fatalf("expected unknown command error for version subcommand, got %v", err)
	}
}
