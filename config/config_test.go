package config

import (
	"reflect"
	"testing"
)

func TestLoadConfig_DefaultValues(t *testing.T) {
	cfg, err := LoadConfig("any-path.yaml")
	if err != nil {
		t.Fatalf("LoadConfig returned unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("LoadConfig returned nil config")
	}

	if cfg.Database.Driver != "postgres" {
		t.Fatalf("unexpected database driver: %q", cfg.Database.Driver)
	}
	if cfg.Database.DSN != "postgres://localhost/db" {
		t.Fatalf("unexpected database DSN: %q", cfg.Database.DSN)
	}
	if cfg.Server.Host != "localhost" {
		t.Fatalf("unexpected server host: %q", cfg.Server.Host)
	}
	if cfg.Server.Port != 8080 {
		t.Fatalf("unexpected server port: %d", cfg.Server.Port)
	}
	if cfg.Extension.Path != "./plugins" {
		t.Fatalf("unexpected extension path: %q", cfg.Extension.Path)
	}
}

func TestLoadConfig_PathIndependence(t *testing.T) {
	first, err := LoadConfig("config/dev.yaml")
	if err != nil {
		t.Fatalf("LoadConfig returned unexpected error for first path: %v", err)
	}

	second, err := LoadConfig("config/prod.yaml")
	if err != nil {
		t.Fatalf("LoadConfig returned unexpected error for second path: %v", err)
	}

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("expected same default config for different paths, got %#v and %#v", first, second)
	}
}
