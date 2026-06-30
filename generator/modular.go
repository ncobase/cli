package generator

import (
	"fmt"

	"github.com/ncobase/cli/generator/templates"
)

func buildModularRenderPlan(data *templates.Data) (*renderPlan, error) {
	registry, err := NewRegistry()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize template registry: %w", err)
	}

	tmplData := NewTemplateData(data)
	plan := newRenderPlan()
	plan.addDir(
		fmt.Sprintf("cmd/%s", data.Name),
		"core",
		"biz",
		"plugin",
		"internal/middleware",
		"internal/server",
		"internal/version",
		"docs",
		"migrations",
		"tests",
	)

	renderRegistryFile := func(name, path string, renderFunc func(*TemplateData) (string, error)) error {
		content, err := renderFunc(tmplData)
		if err != nil {
			return fmt.Errorf("failed to render %s: %w", name, err)
		}
		plan.addFile(path, content)
		return nil
	}
	renderStringFile := func(path, tmpl string) error {
		content, err := renderTemplateString(path, tmpl, data)
		if err != nil {
			return fmt.Errorf("failed to render file %s: %w", path, err)
		}
		plan.addFile(path, content)
		return nil
	}

	if err := renderRegistryFile("main.go", fmt.Sprintf("cmd/%s/main.go", data.Name), registry.RenderMain); err != nil {
		return nil, err
	}
	if err := renderRegistryFile("version.go", "internal/version/version.go", registry.RenderVersion); err != nil {
		return nil, err
	}
	if err := renderRegistryFile("Makefile", "Makefile", registry.RenderMakefile); err != nil {
		return nil, err
	}
	if err := renderRegistryFile(".gitignore", ".gitignore", registry.RenderGitignore); err != nil {
		return nil, err
	}
	if err := renderRegistryFile("config.yaml", "config.yaml", registry.RenderConfigYaml); err != nil {
		return nil, err
	}
	if err := renderStringFile("README.md", templates.ModularReadmeTemplate()); err != nil {
		return nil, err
	}

	plan.addFile("internal/server/server.go", templates.ServerTemplate(data.PackagePath))
	plan.addFile("internal/server/http.go", templates.ModularServerHTTPTemplate(data.PackagePath))
	plan.addFile("internal/server/exts.go", templates.ServerExtsTemplate(data.PackagePath))
	plan.addFile("core/doc.go", templates.PackageDocTemplate("core", "contains foundational product modules"))
	plan.addFile("biz/doc.go", templates.PackageDocTemplate("biz", "contains product business modules"))
	plan.addFile("plugin/doc.go", templates.PackageDocTemplate("plugin", "contains optional integration modules"))
	plan.addFile("docs/README.md", templates.DocsReadmeTemplate())
	plan.addFile("migrations/README.md", templates.MigrationReadmeTemplate())

	if err := renderRegistryFile("CORS middleware", "internal/middleware/cors.go", registry.RenderMiddlewareCORS); err != nil {
		return nil, err
	}
	if err := renderRegistryFile("middleware utils", "internal/middleware/utils.go", registry.RenderMiddlewareUtils); err != nil {
		return nil, err
	}
	plan.addFile("internal/middleware/input.go", templates.MiddlewareInputTemplate())
	if data.WithTracing {
		if err := renderRegistryFile("Trace middleware", "internal/middleware/trace.go", registry.RenderMiddlewareTrace); err != nil {
			return nil, err
		}
	}
	if err := renderRegistryFile("Logger middleware", "internal/middleware/logger.go", registry.RenderMiddlewareLogger); err != nil {
		return nil, err
	}
	if err := renderRegistryFile("Security Headers middleware", "internal/middleware/security_headers.go", registry.RenderMiddlewareSecurityHeaders); err != nil {
		return nil, err
	}
	if err := renderRegistryFile("Client Info middleware", "internal/middleware/client_info.go", registry.RenderMiddlewareClientInfo); err != nil {
		return nil, err
	}

	return plan, nil
}
