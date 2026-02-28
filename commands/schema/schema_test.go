package schema

import "testing"

func TestNewCommand_Metadata(t *testing.T) {
	cmd := NewCommand()

	if cmd.Use != "schema" {
		t.Fatalf("unexpected use field: %q", cmd.Use)
	}
	if cmd.PersistentPreRunE == nil {
		t.Fatal("expected PersistentPreRunE to be configured")
	}

	expectedSubcommands := map[string]bool{
		"inspect": false,
		"apply":   false,
	}

	for _, sub := range cmd.Commands() {
		if _, ok := expectedSubcommands[sub.Name()]; ok {
			expectedSubcommands[sub.Name()] = true
		}
	}

	for name, found := range expectedSubcommands {
		if !found {
			t.Fatalf("expected schema subcommand %q to be registered", name)
		}
	}
}
