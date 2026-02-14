package generator

import (
	"fmt"

	"github.com/ncobase/cli/generator/templates"
)

// Registry manages all templates
type Registry struct {
	loader *Loader
}

// NewRegistry creates a new template registry
func NewRegistry() (*Registry, error) {
	loader, err := NewLoader()
	if err != nil {
		return nil, fmt.Errorf("failed to create template loader: %w", err)
	}

	return &Registry{
		loader: loader,
	}, nil
}

// TemplateData holds common template data
// This should match templates.Data structure
type TemplateData struct {
	Name        string // Extension name
	Type        string // Extension type (core/business/plugin/custom)
	CustomDir   string // Custom directory name, if type is custom
	ModuleName  string // Go module name
	UseMongo    bool   // Whether to use MongoDB
	UseEnt      bool   // Whether to use Ent ORM
	UseGorm     bool   // Whether to use GORM
	WithCmd     bool   // Whether to generate cmd directory with main.go
	WithTest    bool   // Whether to generate test files
	Standalone  bool   // Whether to generate as standalone app
	Group       string // Optional group name
	ExtType     string // Extension type in domain path
	PackagePath string // Full package path

	// Database and data sources
	DBDriver      string
	UseRedis      bool
	UseElastic    bool
	UseOpenSearch bool
	UseMeili      bool
	UseKafka      bool
	UseRabbitMQ   bool
	UseS3Storage  bool
	UseMinio      bool
	UseAliyun     bool
}

// RenderMain renders the main.go template
func (r *Registry) RenderMain(data *TemplateData) (string, error) {
	return r.loader.Render("base/main.go", data)
}

// RenderVersion renders the version.go template
func (r *Registry) RenderVersion(data *TemplateData) (string, error) {
	return r.loader.Render("base/version.go", data)
}

// RenderConfig renders the config.go template
func (r *Registry) RenderConfig(data *TemplateData) (string, error) {
	return r.loader.Render("config/config.go", data)
}

// RenderServer renders the server.go template
func (r *Registry) RenderServer(data *TemplateData) (string, error) {
	return r.loader.Render("server/server.go", data)
}

// RenderHTTP renders the http.go template
func (r *Registry) RenderHTTP(data *TemplateData) (string, error) {
	return r.loader.Render("server/http.go", data)
}

// RenderRest renders the rest.go template
func (r *Registry) RenderRest(data *TemplateData) (string, error) {
	return r.loader.Render("server/rest.go", data)
}

// RenderData renders the appropriate data layer template based on ORM
func (r *Registry) RenderData(data *TemplateData) (string, error) {
	if data.UseEnt {
		return r.loader.Render("data/data-ent.go", data)
	} else if data.UseGorm {
		return r.loader.Render("data/data-gorm.go", data)
	}
	return "", fmt.Errorf("no ORM specified")
}

// RenderModel renders the model.go template
func (r *Registry) RenderModel(data *TemplateData) (string, error) {
	return r.loader.Render("data/model.go", data)
}

// RenderRepositoryProvider renders the repository provider template
func (r *Registry) RenderRepositoryProvider(data *TemplateData) (string, error) {
	return r.loader.Render("data/repository-provider.go", data)
}

// RenderRepository renders the appropriate repository template based on ORM
func (r *Registry) RenderRepository(data *TemplateData) (string, error) {
	if data.UseEnt {
		return r.loader.Render("data/repository-ent.go", data)
	} else if data.UseGorm {
		return r.loader.Render("data/repository-gorm.go", data)
	}
	return "", fmt.Errorf("no ORM specified")
}

// RenderSchema renders the Ent schema template
func (r *Registry) RenderSchema(data *TemplateData) (string, error) {
	if !data.UseEnt {
		return "", fmt.Errorf("Ent ORM not enabled")
	}
	return r.loader.Render("data/schema-user.ent", data)
}

// RenderHandlerProvider renders the handler provider template
func (r *Registry) RenderHandlerProvider(data *TemplateData) (string, error) {
	return r.loader.Render("handler/provider.go", data)
}

// RenderHandler renders the handler template
func (r *Registry) RenderHandler(data *TemplateData) (string, error) {
	return r.loader.Render("handler/handler.go", data)
}

// RenderServiceProvider renders the service provider template
func (r *Registry) RenderServiceProvider(data *TemplateData) (string, error) {
	return r.loader.Render("service/provider.go", data)
}

// RenderService renders the service template
func (r *Registry) RenderService(data *TemplateData) (string, error) {
	return r.loader.Render("service/service.go", data)
}

// RenderMakefile renders the Makefile template
func (r *Registry) RenderMakefile(data *TemplateData) (string, error) {
	return r.loader.Render("base/makefile", data)
}

// RenderGitignore renders the .gitignore template
func (r *Registry) RenderGitignore(data *TemplateData) (string, error) {
	return r.loader.Render("base/gitignore", data)
}

// RenderReadme renders the README.md template
func (r *Registry) RenderReadme(data *TemplateData) (string, error) {
	return r.loader.Render("base/readme", data)
}

// RenderConfigYaml renders the config.yaml template
func (r *Registry) RenderConfigYaml(data *TemplateData) (string, error) {
	return r.loader.Render("base/config.yaml", data)
}

// RenderGenerate renders the generate.go template for Ent
func (r *Registry) RenderGenerate(data *TemplateData) (string, error) {
	return r.loader.Render("base/generate.go", data)
}

// RenderMiddlewareCORS renders the CORS middleware template
func (r *Registry) RenderMiddlewareCORS(data *TemplateData) (string, error) {
	return r.loader.Render("middleware/cors.go", data)
}

// RenderMiddlewareTrace renders the Trace middleware template
func (r *Registry) RenderMiddlewareTrace(data *TemplateData) (string, error) {
	return r.loader.Render("middleware/trace.go", data)
}

// RenderMiddlewareLogger renders the Logger middleware template
func (r *Registry) RenderMiddlewareLogger(data *TemplateData) (string, error) {
	return r.loader.Render("middleware/logger.go", data)
}

// RenderMiddlewareSecurityHeaders renders the Security Headers middleware template
func (r *Registry) RenderMiddlewareSecurityHeaders(data *TemplateData) (string, error) {
	return r.loader.Render("middleware/security_headers.go", data)
}

// RenderMiddlewareClientInfo renders the Client Info middleware template
func (r *Registry) RenderMiddlewareClientInfo(data *TemplateData) (string, error) {
	return r.loader.Render("middleware/client_info.go", data)
}

// RenderMiddlewareUtils renders the middleware utilities template
func (r *Registry) RenderMiddlewareUtils(data *TemplateData) (string, error) {
	return r.loader.Render("middleware/utils.go", data)
}

// NewTemplateData creates template data from templates.Data
func NewTemplateData(d *templates.Data) *TemplateData {
	return &TemplateData{
		Name:          d.Name,
		Type:          d.Type,
		CustomDir:     d.CustomDir,
		ModuleName:    d.ModuleName,
		UseMongo:      d.UseMongo,
		UseEnt:        d.UseEnt,
		UseGorm:       d.UseGorm,
		WithCmd:       d.WithCmd,
		WithTest:      d.WithTest,
		Standalone:    d.Standalone,
		Group:         d.Group,
		ExtType:       d.ExtType,
		PackagePath:   d.PackagePath,
		DBDriver:      d.DBDriver,
		UseRedis:      d.UseRedis,
		UseElastic:    d.UseElastic,
		UseOpenSearch: d.UseOpenSearch,
		UseMeili:      d.UseMeili,
		UseKafka:      d.UseKafka,
		UseRabbitMQ:   d.UseRabbitMQ,
		UseS3Storage:  d.UseS3Storage,
		UseMinio:      d.UseMinio,
		UseAliyun:     d.UseAliyun,
	}
}
