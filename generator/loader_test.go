package generator

import (
	"strings"
	"testing"
)

func TestLoaderRender(t *testing.T) {
	loader, err := NewLoader()
	if err != nil {
		t.Fatalf("NewLoader() error = %v", err)
	}

	out, err := loader.Render("base/gitignore", nil)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if !strings.Contains(out, ".env") {
		t.Fatalf("expected rendered gitignore to contain .env, got: %q", out)
	}
}

func TestLoaderRender_NotFound(t *testing.T) {
	loader, err := NewLoader()
	if err != nil {
		t.Fatalf("NewLoader() error = %v", err)
	}

	_, err = loader.Render("unknown/template", nil)
	if err == nil {
		t.Fatal("expected error for missing template, got nil")
	}
}

func TestRegistryRenderData(t *testing.T) {
	registry, err := NewRegistry()
	if err != nil {
		t.Fatalf("NewRegistry() error = %v", err)
	}

	_, err = registry.RenderData(&TemplateData{})
	if err == nil {
		t.Fatal("expected error when no ORM is enabled")
	}

	out, err := registry.RenderData(&TemplateData{UseEnt: true})
	if err != nil {
		t.Fatalf("RenderData() with Ent error = %v", err)
	}
	if strings.TrimSpace(out) == "" {
		t.Fatal("expected non-empty rendered data template")
	}
}
