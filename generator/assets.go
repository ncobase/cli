package generator

import (
	"fmt"

	"github.com/ncobase/cli/generator/templates"
)

func buildStandaloneRenderPlan(data *templates.Data) (*renderPlan, error) {
	registry, err := NewRegistry()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize template registry: %w", err)
	}

	tmplData := NewTemplateData(data)
	plan := newRenderPlan()

	directories := []string{
		fmt.Sprintf("cmd/%s", data.Name),
		"internal/server", "internal/middleware", "internal/version", "internal/config",
		"handler", "data", "data/model", "data/repository", "service",
	}

	if data.UseEnt {
		directories = append(directories, "data/ent", "data/schema")
	}
	directories = append(directories, "data/migrates")
	if data.WithTest {
		directories = append(directories, "tests")
	}
	plan.addDir(directories...)

	renderFile := func(name, path string, renderFunc func(*TemplateData) (string, error)) error {
		content, err := renderFunc(tmplData)
		if err != nil {
			return fmt.Errorf("failed to render %s: %w", name, err)
		}
		plan.addFile(path, content)
		return nil
	}

	if err := renderFile("main.go", fmt.Sprintf("cmd/%s/main.go", data.Name), registry.RenderMain); err != nil {
		return nil, err
	}
	if err := renderFile("version.go", "internal/version/version.go", registry.RenderVersion); err != nil {
		return nil, err
	}
	if err := renderFile("Makefile", "Makefile", registry.RenderMakefile); err != nil {
		return nil, err
	}
	if err := renderFile(".gitignore", ".gitignore", registry.RenderGitignore); err != nil {
		return nil, err
	}
	if err := renderFile("README.md", "README.md", registry.RenderReadme); err != nil {
		return nil, err
	}
	if err := renderFile("config.yaml", "config.yaml", registry.RenderConfigYaml); err != nil {
		return nil, err
	}
	if err := renderFile("server.go", "internal/server/server.go", registry.RenderServer); err != nil {
		return nil, err
	}
	if err := renderFile("exts.go", "internal/server/exts.go", registry.RenderServerExts); err != nil {
		return nil, err
	}
	if err := renderFile("http.go", "internal/server/http.go", registry.RenderHTTP); err != nil {
		return nil, err
	}
	if err := renderFile("rest.go", "internal/server/rest.go", registry.RenderRest); err != nil {
		return nil, err
	}

	if data.WithGRPC {
		if err := renderFile("grpc.go", "internal/server/grpc.go", registry.RenderGRPCServer); err != nil {
			return nil, err
		}
	}

	if err := renderFile("config.go", "internal/config/config.go", registry.RenderConfig); err != nil {
		return nil, err
	}
	if err := renderFile("handler provider.go", "handler/provider.go", registry.RenderHandlerProvider); err != nil {
		return nil, err
	}
	if err := renderFile("handler.go", "handler/handler.go", registry.RenderHandler); err != nil {
		return nil, err
	}
	if err := renderFile("service provider.go", "service/provider.go", registry.RenderServiceProvider); err != nil {
		return nil, err
	}
	if err := renderFile("service.go", "service/service.go", registry.RenderService); err != nil {
		return nil, err
	}
	if err := renderFile("data.go", "data/data.go", registry.RenderData); err != nil {
		return nil, err
	}
	if err := renderFile("model.go", "data/model/model.go", registry.RenderModel); err != nil {
		return nil, err
	}
	if err := renderFile("repository provider.go", "data/repository/provider.go", registry.RenderRepositoryProvider); err != nil {
		return nil, err
	}
	if err := renderFile("repository.go", "data/repository/repository.go", registry.RenderRepository); err != nil {
		return nil, err
	}
	if err := renderFile("CORS middleware", "internal/middleware/cors.go", registry.RenderMiddlewareCORS); err != nil {
		return nil, err
	}
	if err := renderFile("middleware utils", "internal/middleware/utils.go", registry.RenderMiddlewareUtils); err != nil {
		return nil, err
	}

	if data.WithTracing {
		if err := renderFile("Trace middleware", "internal/middleware/trace.go", registry.RenderMiddlewareTrace); err != nil {
			return nil, err
		}
	}

	if err := renderFile("Logger middleware", "internal/middleware/logger.go", registry.RenderMiddlewareLogger); err != nil {
		return nil, err
	}
	if err := renderFile("Security Headers middleware", "internal/middleware/security_headers.go", registry.RenderMiddlewareSecurityHeaders); err != nil {
		return nil, err
	}
	if err := renderFile("Client Info middleware", "internal/middleware/client_info.go", registry.RenderMiddlewareClientInfo); err != nil {
		return nil, err
	}

	if data.WithTest {
		if err := renderFile("handler test", "tests/handler_test.go", registry.RenderHandlerTest); err != nil {
			return nil, err
		}
		if err := renderFile("service test", "tests/service_test.go", registry.RenderServiceTest); err != nil {
			return nil, err
		}
	}

	if data.UseEnt {
		if err := renderFile("schema", "data/schema/example.go", registry.RenderSchema); err != nil {
			return nil, err
		}
		if err := renderFile("generate.go", "generate.go", registry.RenderGenerate); err != nil {
			return nil, err
		}
	}

	return plan, nil
}

// getDesc returns the description of the generated component
func getDesc(data *templates.Data) string {
	if data.Standalone && data.ProjectType == "modular" {
		return "modular application"
	}
	if data.Standalone {
		return "standalone application"
	}
	if data.Type == "custom" {
		return fmt.Sprintf("'%s' directory", data.CustomDir)
	}
	if data.Type == "direct" {
		return fmt.Sprintf("'%s' directory", data.Name)
	}
	extDescriptions := map[string]string{
		"core":     "Core Domain",
		"biz":      "Business Domain",
		"business": "Business Domain",
		"plugin":   "Plugin Domain",
		"custom":   "Custom Directory",
	}
	return extDescriptions[data.ExtType]
}
