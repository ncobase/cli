// Package generator provides code generation capabilities for ncobase CLI.
//
// It supports two generation modes:
//
// 1. Standalone Application Mode (via 'nco init'):
//   - Creates a complete, ready-to-run Go application
//   - Generates full project structure: cmd/, data/, handler/, service/, etc.
//   - Includes configuration files, Makefile, and documentation
//   - Supports multiple databases, ORMs, and data sources
//
// 2. Extension Module Mode (via 'nco create'):
//   - Creates extension modules within existing ncobase projects
//   - Supports three extension types: core, business, plugin
//   - Allows custom directory locations
//   - Integrates with existing project structure
//
// The generator uses an embed.FS-based template system (see loader.go)
// for efficient template management and rendering.
package generator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ncobase/cli/generator/templates"
	"github.com/ncobase/cli/utils"
)

// Options defines code generation options
type Options struct {
	Name       string // Project or extension name
	Type       string // Generation type: core, business, plugin, custom, or direct
	CustomDir  string // Custom directory name (when Type is custom)
	OutputPath string // Base output directory
	ModuleName string // Go module name (e.g., github.com/username/project)

	// ORM options (mutually exclusive)
	UseMongo bool // Use MongoDB driver
	UseEnt   bool // Use Ent ORM for SQL databases
	UseGorm  bool // Use GORM for SQL databases

	// Generation options
	WithCmd  bool   // Generate cmd directory with main.go
	WithTest bool   // Generate test files (unit, integration, e2e)
	Group    string // Optional domain group name

	// Standalone mode (set by init command)
	Standalone bool // Generate as standalone application

	// Database configuration
	DBDriver string // Database driver: postgres, mysql, sqlite, mongodb, neo4j

	// Data source drivers
	UseRedis      bool // Include Redis driver for caching/queuing
	UseElastic    bool // Include Elasticsearch driver for search
	UseOpenSearch bool // Include OpenSearch driver for search
	UseMeili      bool // Include Meilisearch driver for search

	// Message queue drivers
	UseKafka    bool // Include Kafka driver for messaging
	UseRabbitMQ bool // Include RabbitMQ driver for messaging

	// Storage drivers
	UseS3Storage bool // Include AWS S3 storage driver
	UseMinio     bool // Include MinIO storage driver
	UseAliyun    bool // Include Aliyun OSS storage driver
}

// DefaultOptions returns default options
func DefaultOptions() *Options {
	return &Options{
		Type:          "custom",
		OutputPath:    "",
		ModuleName:    "",
		UseMongo:      false,
		UseEnt:        false,
		UseGorm:       false,
		WithCmd:       false,
		WithTest:      false,
		Standalone:    false,
		Group:         "",
		DBDriver:      "",
		UseRedis:      false,
		UseElastic:    false,
		UseOpenSearch: false,
		UseMeili:      false,
		UseKafka:      false,
		UseRabbitMQ:   false,
		UseS3Storage:  false,
		UseMinio:      false,
		UseAliyun:     false,
	}
}

// extDescriptions maps extension types to human-readable descriptions
var extDescriptions = map[string]string{
	"core":     "Core Domain",      // Fundamental business logic
	"business": "Business Domain",  // Application-specific logic
	"plugin":   "Plugin Domain",    // Optional features
	"custom":   "Custom Directory", // User-defined location
}

// Generate generates code
func Generate(opts *Options) error {
	if !utils.ValidateName(opts.Name) {
		return fmt.Errorf("invalid name: %s", opts.Name)
	}

	if opts.UseMongo && opts.DBDriver == "" {
		opts.DBDriver = "mongodb"
	}
	if opts.DBDriver == "" && (opts.UseEnt || opts.UseGorm) {
		opts.DBDriver = "all"
	}

	// Determine output path
	if opts.OutputPath == "" {
		// Use current directory
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %v", err)
		}
		opts.OutputPath = cwd
	}

	// Determine module name if not provided
	if opts.ModuleName == "" {
		if opts.Standalone {
			// For standalone apps, default module name is the project name
			opts.ModuleName = opts.Name
		} else {
			// Try to detect from go.mod file
			goModPath := filepath.Join(opts.OutputPath, "go.mod")
			if utils.FileExists(goModPath) {
				content, err := os.ReadFile(goModPath)
				if err == nil {
					lines := strings.Split(string(content), "\n")
					for _, line := range lines {
						if strings.HasPrefix(line, "module ") {
							opts.ModuleName = strings.TrimSpace(strings.TrimPrefix(line, "module "))
							break
						}
					}
				}
			}

			// If still empty, use a default
			if opts.ModuleName == "" {
				// Use current directory name as module name
				dirs := strings.Split(opts.OutputPath, string(os.PathSeparator))
				opts.ModuleName = dirs[len(dirs)-1]
			}
		}
	}

	var basePath string
	var extType string
	var mainTemplate func(string) string

	// Handle standalone mode differently
	if opts.Standalone {
		// In standalone mode, generate only the cmd directory
		if opts.Type == "direct" {
			basePath = filepath.Join(opts.OutputPath, opts.Name)
		} else if opts.Type == "custom" {
			basePath = filepath.Join(opts.OutputPath, opts.CustomDir, opts.Name)
		} else {
			basePath = filepath.Join(opts.OutputPath, opts.Type, opts.Name)
		}

		extType = opts.Type

		// Check if directory already exists
		if utils.PathExists(basePath) {
			return fmt.Errorf("directory '%s' already exists", basePath)
		}

		// Create base directory
		if err := utils.EnsureDir(basePath); err != nil {
			return fmt.Errorf("failed to create base directory: %v", err)
		}

		// Prepare template data
		data := &templates.Data{
			Name:          opts.Name,
			Type:          opts.Type,
			UseMongo:      opts.UseMongo,
			UseEnt:        opts.UseEnt,
			UseGorm:       opts.UseGorm,
			WithTest:      opts.WithTest,
			WithCmd:       true, // Standalone always includes cmd
			Standalone:    opts.Standalone,
			Group:         opts.Group,
			ExtType:       extType,
			ModuleName:    opts.ModuleName,
			CustomDir:     opts.CustomDir,
			PackagePath:   getPackagePath(opts),
			DBDriver:      opts.DBDriver,
			UseRedis:      opts.UseRedis,
			UseElastic:    opts.UseElastic,
			UseOpenSearch: opts.UseOpenSearch,
			UseMeili:      opts.UseMeili,
			UseKafka:      opts.UseKafka,
			UseRabbitMQ:   opts.UseRabbitMQ,
			UseS3Storage:  opts.UseS3Storage,
			UseMinio:      opts.UseMinio,
			UseAliyun:     opts.UseAliyun,
		}

		// Create standalone structure
		if err := createStandaloneStructure(basePath, data); err != nil {
			return err
		}

		// Initialize Go module for standalone mode
		if err := initializeGoModule(basePath, data, opts); err != nil {
			fmt.Printf("Warning: failed to initialize Go module: %v\n", err)
			// Don't interrupt the flow, just warn
		}

		fmt.Printf("Successfully generated standalone application '%s' in %s\n", data.Name, getDesc(data))
		return nil
	}

	// Regular extension generation (not standalone)
	// Determine base paths and templates based on type
	switch opts.Type {
	case "core":
		basePath = filepath.Join(opts.OutputPath, "core", opts.Name)
		extType = "core"
		mainTemplate = templates.CoreTemplate
	case "business":
		basePath = filepath.Join(opts.OutputPath, "business", opts.Name)
		extType = "business"
		mainTemplate = templates.BusinessTemplate
	case "plugin":
		basePath = filepath.Join(opts.OutputPath, "plugin", opts.Name)
		extType = "plugin"
		mainTemplate = templates.PluginTemplate
	case "direct":
		basePath = filepath.Join(opts.OutputPath, opts.Name)
		extType = "direct"
		// Use business template
		mainTemplate = templates.BusinessTemplate
	case "custom":
		basePath = filepath.Join(opts.OutputPath, opts.CustomDir, opts.Name)
		extType = "custom"
		// Use business template
		mainTemplate = templates.BusinessTemplate
	default:
		return fmt.Errorf("unknown type: %s", opts.Type)
	}

	// Check if component already exists
	if utils.PathExists(basePath) {
		return fmt.Errorf("'%s' already exists in %s", opts.Name, extDescriptions[extType])
	}

	// Prepare template data
	data := &templates.Data{
		Name:          opts.Name,
		Type:          opts.Type,
		UseMongo:      opts.UseMongo,
		UseEnt:        opts.UseEnt,
		UseGorm:       opts.UseGorm,
		WithTest:      opts.WithTest,
		WithCmd:       opts.WithCmd,
		Standalone:    opts.Standalone,
		Group:         opts.Group,
		ExtType:       extType,
		ModuleName:    opts.ModuleName,
		CustomDir:     opts.CustomDir,
		PackagePath:   getPackagePath(opts),
		DBDriver:      opts.DBDriver,
		UseRedis:      opts.UseRedis,
		UseElastic:    opts.UseElastic,
		UseOpenSearch: opts.UseOpenSearch,
		UseMeili:      opts.UseMeili,
		UseKafka:      opts.UseKafka,
		UseRabbitMQ:   opts.UseRabbitMQ,
		UseS3Storage:  opts.UseS3Storage,
		UseMinio:      opts.UseMinio,
		UseAliyun:     opts.UseAliyun,
	}

	// Create the main extension structure
	err := createStructure(basePath, data, mainTemplate)
	if err != nil {
		return err
	}

	// Generate cmd directory if WithCmd is true
	if opts.WithCmd {
		// Create cmd directory inside the extension directory
		cmdPath := filepath.Join(basePath, "cmd")
		if err := utils.EnsureDir(cmdPath); err != nil {
			return fmt.Errorf("failed to create cmd directory: %v", err)
		}

		// Create internal/server directory
		serverPath := filepath.Join(basePath, "internal/server")
		if err := utils.EnsureDir(serverPath); err != nil {
			return fmt.Errorf("failed to create internal/server directory: %v", err)
		}

		// Create internal/middleware directory
		middlewarePath := filepath.Join(basePath, "internal/middleware")
		if err := utils.EnsureDir(middlewarePath); err != nil {
			return fmt.Errorf("failed to create internal/middleware directory: %v", err)
		}

		// Create internal/version directory
		versionPath := filepath.Join(basePath, "internal/version")
		if err := utils.EnsureDir(versionPath); err != nil {
			return fmt.Errorf("failed to create internal/version directory: %v", err)
		}

		// Create files in cmd directory
		files := map[string]string{
			"cmd/main.go": templates.CmdMainTemplate(data),

			// Internal Server
			"internal/server/server.go": templates.ServerTemplate(data.PackagePath),
			"internal/server/http.go":   templates.ServerHTTPTemplate(data.PackagePath),
			"internal/server/exts.go":   templates.ServerExtsTemplate(data.PackagePath),

			// Internal Middleware
			"internal/middleware/cors.go":             templates.MiddlewareCORSTemplate(),
			"internal/middleware/security_headers.go": templates.MiddlewareSecurityHeadersTemplate(),
			"internal/middleware/trace.go":            templates.MiddlewareTraceTemplate(),
			"internal/middleware/logger.go":           templates.MiddlewareLoggerTemplate(),
			"internal/middleware/client_info.go":      templates.MiddlewareClientInfoTemplate(),

			// Internal Version
			"internal/version/version.go": templates.VersionTemplate(),
		}

		// Write files
		for filePath, tmpl := range files {
			if err := utils.WriteTemplateFile(
				filepath.Join(basePath, filePath),
				tmpl,
				data,
			); err != nil {
				return fmt.Errorf("failed to create file %s: %v", filePath, err)
			}
		}

		// Initialize Go module for WithCmd mode
		if err := initializeGoModule(basePath, data, opts); err != nil {
			fmt.Printf("Warning: failed to initialize Go module: %v\n", err)
			// Don't interrupt the flow, just warn
		}
	}

	fmt.Printf("Successfully generated '%s' in %s\n", data.Name, getDesc(data))
	return nil
}

// getPackagePath returns the package path based on options
func getPackagePath(opts *Options) string {
	if opts.Standalone {
		return opts.ModuleName
	}
	switch opts.Type {
	case "custom":
		if opts.CustomDir == "" {
			return fmt.Sprintf("%s/%s", opts.ModuleName, opts.Name)
		}
		return fmt.Sprintf("%s/%s/%s", opts.ModuleName, opts.CustomDir, opts.Name)
	case "direct":
		return fmt.Sprintf("%s/%s", opts.ModuleName, opts.Name)
	default:
		return fmt.Sprintf("%s/%s/%s", opts.ModuleName, opts.Type, opts.Name)
	}
}

func createStructure(basePath string, data *templates.Data, mainTemplate func(string) string) error {
	// Create base directory
	if err := utils.EnsureDir(basePath); err != nil {
		return fmt.Errorf("failed to create base directory: %v", err)
	}

	// Create directory structure
	directories := []string{
		"data",
		"data/repository",
		"data/schema",
		"handler",
		"service",
		"structs",
	}

	if data.WithTest {
		directories = append(directories, "tests")
	}

	for _, dir := range directories {
		if err := utils.EnsureDir(filepath.Join(basePath, dir)); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	// Create files
	selectDataTemplate := func(data templates.Data) string {
		if data.UseEnt {
			return templates.DataTemplateWithEnt(data.Name, data.ExtType)
		}
		if data.UseGorm {
			return templates.DataTemplateWithGorm(data.Name, data.ExtType)
		}
		if data.UseMongo {
			return templates.DataTemplateWithMongo(data.Name, data.ExtType)
		}
		return templates.DataTemplate(data.Name, data.ExtType)
	}

	files := map[string]string{
		fmt.Sprintf("%s.go", data.Name): mainTemplate(data.Name),
		"data/data.go":                  selectDataTemplate(*data),
		"data/repository/provider.go":   templates.RepositoryTemplate(data),
		"data/schema/schema.go":         templates.SchemaTemplate(),
		"handler/provider.go":           templates.HandlerTemplate(data.Name, data.ExtType, data.ModuleName),
		"service/provider.go":           templates.ServiceTemplate(data.Name, data.ExtType, data.ModuleName),
		"structs/structs.go":            templates.StructsTemplate(),
	}

	// Add ent files if required
	if data.UseEnt {
		files["generate.go"] = templates.GeneraterTemplate(data.Name, data.ExtType, data.ModuleName)
	}

	// Add test files if required
	if data.WithTest {
		files["tests/ext_test.go"] = templates.ExtTestTemplate(data.Name, data.ExtType, data.ModuleName)
		files["tests/handler_test.go"] = templates.HandlerTestTemplate(data.Name, data.ExtType, data.ModuleName)
		files["tests/service_test.go"] = templates.ServiceTestTemplate(data.Name, data.ExtType, data.ModuleName)
	}

	// Write files
	for filePath, tmpl := range files {
		if err := utils.WriteTemplateFile(
			filepath.Join(basePath, filePath),
			tmpl,
			data,
		); err != nil {
			return fmt.Errorf("failed to create file %s: %v", filePath, err)
		}
	}

	return nil
}

// createStandaloneStructure creates the structure for a standalone application
func createStandaloneStructure(basePath string, data *templates.Data) error {
	// Initialize template registry
	registry, err := NewRegistry()
	if err != nil {
		return fmt.Errorf("failed to initialize template registry: %w", err)
	}

	// Create template data
	tmplData := NewTemplateData(data)

	// Create essential directories
	directories := []string{
		fmt.Sprintf("cmd/%s", data.Name),
		"internal/server",
		"internal/middleware",
		"internal/version",
		"internal/config",
		"handler",
		"data",
		"data/model",
		"data/repository",
		"service",
	}

	if data.UseEnt {
		directories = append(directories, "data/ent", "data/schema")
	}

	// Add migrates directory
	directories = append(directories, "data/migrates")

	if data.WithTest {
		directories = append(directories, "tests")
	}

	for _, dir := range directories {
		if err := utils.EnsureDir(filepath.Join(basePath, dir)); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	// Prepare file content map
	files := make(map[string]string)

	// Generate base files
	if content, err := registry.RenderMain(tmplData); err == nil {
		files[fmt.Sprintf("cmd/%s/main.go", data.Name)] = content
	} else {
		return fmt.Errorf("failed to render main.go: %w", err)
	}

	if content, err := registry.RenderVersion(tmplData); err == nil {
		files["internal/version/version.go"] = content
	} else {
		return fmt.Errorf("failed to render version.go: %w", err)
	}

	if content, err := registry.RenderMakefile(tmplData); err == nil {
		files["Makefile"] = content
	} else {
		return fmt.Errorf("failed to render Makefile: %w", err)
	}

	if content, err := registry.RenderGitignore(tmplData); err == nil {
		files[".gitignore"] = content
	} else {
		return fmt.Errorf("failed to render .gitignore: %w", err)
	}

	if content, err := registry.RenderReadme(tmplData); err == nil {
		files["README.md"] = content
	} else {
		return fmt.Errorf("failed to render README.md: %w", err)
	}

	if content, err := registry.RenderConfigYaml(tmplData); err == nil {
		files["config.yaml"] = content
	} else {
		return fmt.Errorf("failed to render config.yaml: %w", err)
	}

	// Server layer
	if content, err := registry.RenderServer(tmplData); err == nil {
		files["internal/server/server.go"] = content
	} else {
		return fmt.Errorf("failed to render server.go: %w", err)
	}

	if content, err := registry.RenderHTTP(tmplData); err == nil {
		files["internal/server/http.go"] = content
	} else {
		return fmt.Errorf("failed to render http.go: %w", err)
	}

	if content, err := registry.RenderRest(tmplData); err == nil {
		files["internal/server/rest.go"] = content
	} else {
		return fmt.Errorf("failed to render rest.go: %w", err)
	}

	// Config layer
	if content, err := registry.RenderConfig(tmplData); err == nil {
		files["internal/config/config.go"] = content
	} else {
		return fmt.Errorf("failed to render config.go: %w", err)
	}

	// Handler layer
	if content, err := registry.RenderHandlerProvider(tmplData); err == nil {
		files["handler/provider.go"] = content
	} else {
		return fmt.Errorf("failed to render handler provider.go: %w", err)
	}

	if content, err := registry.RenderHandler(tmplData); err == nil {
		files["handler/handler.go"] = content
	} else {
		return fmt.Errorf("failed to render handler.go: %w", err)
	}

	// Service layer
	if content, err := registry.RenderServiceProvider(tmplData); err == nil {
		files["service/provider.go"] = content
	} else {
		return fmt.Errorf("failed to render service provider.go: %w", err)
	}

	if content, err := registry.RenderService(tmplData); err == nil {
		files["service/service.go"] = content
	} else {
		return fmt.Errorf("failed to render service.go: %w", err)
	}

	// Data layer
	if content, err := registry.RenderData(tmplData); err == nil {
		files["data/data.go"] = content
	} else {
		return fmt.Errorf("failed to render data.go: %w", err)
	}

	if content, err := registry.RenderModel(tmplData); err == nil {
		files["data/model/model.go"] = content
	} else {
		return fmt.Errorf("failed to render model.go: %w", err)
	}

	// Repository layer
	if content, err := registry.RenderRepositoryProvider(tmplData); err == nil {
		files["data/repository/provider.go"] = content
	} else {
		return fmt.Errorf("failed to render repository provider.go: %w", err)
	}

	if content, err := registry.RenderRepository(tmplData); err == nil {
		files["data/repository/repository.go"] = content
	} else {
		return fmt.Errorf("failed to render repository.go: %w", err)
	}

	// Middleware layer
	if content, err := registry.RenderMiddlewareCORS(tmplData); err == nil {
		files["internal/middleware/cors.go"] = content
	} else {
		return fmt.Errorf("failed to render CORS middleware: %w", err)
	}

	if content, err := registry.RenderMiddlewareTrace(tmplData); err == nil {
		files["internal/middleware/trace.go"] = content
	} else {
		return fmt.Errorf("failed to render Trace middleware: %w", err)
	}

	if content, err := registry.RenderMiddlewareLogger(tmplData); err == nil {
		files["internal/middleware/logger.go"] = content
	} else {
		return fmt.Errorf("failed to render Logger middleware: %w", err)
	}

	if content, err := registry.RenderMiddlewareSecurityHeaders(tmplData); err == nil {
		files["internal/middleware/security_headers.go"] = content
	} else {
		return fmt.Errorf("failed to render Security Headers middleware: %w", err)
	}

	if content, err := registry.RenderMiddlewareClientInfo(tmplData); err == nil {
		files["internal/middleware/client_info.go"] = content
	} else {
		return fmt.Errorf("failed to render Client Info middleware: %w", err)
	}

	if content, err := registry.RenderMiddlewareUtils(tmplData); err == nil {
		files["internal/middleware/utils.go"] = content
	} else {
		return fmt.Errorf("failed to render middleware utils: %w", err)
	}

	// Add test files if required
	if data.WithTest {
		files["tests/handler_test.go"] = templates.StandaloneHandlerTestTemplate(data.Name, data.ModuleName)
		files["tests/service_test.go"] = templates.StandaloneServiceTestTemplate(data.Name, data.ModuleName)
	}

	// Add schema if Ent
	if data.UseEnt {
		if content, err := registry.RenderSchema(tmplData); err == nil {
			files["data/schema/user.go"] = content
		} else {
			return fmt.Errorf("failed to render schema: %w", err)
		}
		if content, err := registry.RenderGenerate(tmplData); err == nil {
			files["generate.go"] = content
		} else {
			return fmt.Errorf("failed to render generate.go: %w", err)
		}
	}

	// Write all files
	for filePath, content := range files {
		fullPath := filepath.Join(basePath, filePath)
		if err := utils.EnsureDir(filepath.Dir(fullPath)); err != nil {
			return fmt.Errorf("failed to create directory for %s: %v", filePath, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %v", filePath, err)
		}
	}

	return nil
}

// getDesc returns the description of the generated component
func getDesc(data *templates.Data) string {
	if data.Type == "custom" {
		return fmt.Sprintf("'%s' directory", data.CustomDir)
	}
	if data.Type == "direct" {
		return fmt.Sprintf("'%s' directory", data.Name)
	}
	return extDescriptions[data.ExtType]
}

// initializeGoModule initializes a Go module for the generated code
// This is used for both standalone and with-cmd modes
func initializeGoModule(basePath string, data *templates.Data, opts *Options) error {
	// Create go.mod file
	goModPath := filepath.Join(basePath, "go.mod")

	// Use strings.Builder for efficient string concatenation
	var builder strings.Builder

	// Create initial go.mod content
	fmt.Fprintf(&builder, `module %s

go 1.25.5

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/google/uuid v1.6.0
)
`, data.PackagePath)

	// Add ncore dependencies with versions
	ncoreDeps := []string{
		"github.com/ncobase/ncore/config",
		"github.com/ncobase/ncore/logging",
		"github.com/ncobase/ncore/ecode",
		"github.com/ncobase/ncore/net",
		"github.com/ncobase/ncore/extension",
	}

	// Add database-specific dependencies
	if opts.UseMongo {
		builder.WriteString("\nrequire go.mongodb.org/mongo-driver v1.17.6\n")
	}

	if opts.UseEnt {
		builder.WriteString("\nrequire entgo.io/ent v0.14.1\n")
	}

	if opts.UseGorm {
		builder.WriteString(`
require (
	gorm.io/gorm v1.25.12
	gorm.io/driver/mysql v1.5.7
	gorm.io/driver/postgres v1.5.11
	gorm.io/driver/sqlite v1.5.7
)
`)
	}

	// Add data driver dependencies
	if opts.DBDriver != "" && opts.DBDriver != "none" {
		ncoreDeps = append(ncoreDeps, "github.com/ncobase/ncore/data/"+opts.DBDriver)
	}

	if opts.UseRedis {
		ncoreDeps = append(ncoreDeps, "github.com/ncobase/ncore/data/redis")
	}

	if opts.UseElastic {
		ncoreDeps = append(ncoreDeps, "github.com/ncobase/ncore/data/elasticsearch")
	}

	if opts.UseOpenSearch {
		ncoreDeps = append(ncoreDeps, "github.com/ncobase/ncore/data/opensearch")
	}

	if opts.UseMeili {
		ncoreDeps = append(ncoreDeps, "github.com/ncobase/ncore/data/meilisearch")
	}

	if opts.UseKafka {
		ncoreDeps = append(ncoreDeps, "github.com/ncobase/ncore/data/kafka")
	}

	if opts.UseRabbitMQ {
		ncoreDeps = append(ncoreDeps, "github.com/ncobase/ncore/data/rabbitmq")
	}

	if opts.UseS3Storage {
		ncoreDeps = append(ncoreDeps, "github.com/ncobase/ncore/data/s3")
	}

	if opts.UseMinio {
		ncoreDeps = append(ncoreDeps, "github.com/ncobase/ncore/data/minio")
	}

	if opts.UseAliyun {
		ncoreDeps = append(ncoreDeps, "github.com/ncobase/ncore/data/aliyun")
	}

	// Add replace directives for local development
	builder.WriteString("\n// Replace directives for local ncore development\n")
	builder.WriteString("// Remove these lines and run 'go mod tidy' when ncore packages are published\n")
	for _, dep := range ncoreDeps {
		// Extract module name from path (e.g., "config" from "github.com/ncobase/ncore/config")
		parts := strings.Split(dep, "/")
		moduleName := parts[len(parts)-1]
		fmt.Fprintf(&builder, "// replace %s => ../../ncore/%s\n", dep, moduleName)
	}

	// Write go.mod file
	if err := utils.WriteTemplateFile(goModPath, builder.String(), nil); err != nil {
		return fmt.Errorf("failed to create go.mod file: %v", err)
	}

	// Create .gitignore file
	gitignorePath := filepath.Join(basePath, ".gitignore")
	gitignoreContent := `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool
*.out

# Dependency directories
vendor/

# IDE files
.idea/
.vscode/
*.sublime-workspace

# OS specific files
.DS_Store
Thumbs.db
`

	if err := utils.WriteTemplateFile(gitignorePath, gitignoreContent, nil); err != nil {
		fmt.Printf("Warning: failed to create .gitignore file: %v\n", err)
		// Just warn, don't stop the process
	}

	// Execute go mod tidy to resolve dependencies
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = basePath
	if err := tidyCmd.Run(); err != nil {
		fmt.Printf("Warning: failed to run 'go mod tidy': %v\n", err)
		// Just warn, don't stop the process
	}

	// Initialize additional tools based on options
	if opts.UseEnt {
		// Ensure schema directory exists
		schemaDir := filepath.Join(basePath, "data/schema")
		if err := utils.EnsureDir(schemaDir); err != nil {
			fmt.Printf("Warning: failed to create ent schema directory: %v\n", err)
			return nil
		}
	}

	// Create a basic README.md
	readmePath := filepath.Join(basePath, "README.md")
	readmeContent := fmt.Sprintf(`# %s

## Overview

This application was generated using the Ncobase CLI. It follows the standard Ncobase architecture and best practices, employing a clean architecture design with domain-driven principles.

## Project Structure

### Core Directories

- **cmd/**: Application entry points.
  - **main.go**: The main entry point that initializes the application.
- **internal/**: Private application code.
  - **server/**: Server initialization, HTTP configuration, and extension management.
  - **middleware/**: Application-specific HTTP middleware (CORS, Logger, Tracing, etc.).
  - **version/**: Version information handling.
- **config/**: Configuration management and structure definitions.
- **handler/**: HTTP handlers (Controllers) responsible for processing requests and returning responses.
- **service/**: Business logic layer where the core application logic resides.
- **data/**: Data access layer.
  - **model/**: Domain models and data entities.
  - **repository/**: Database interaction logic (Repository pattern).
  - **ent/**: (Optional) Entity framework generated code.
  - **schema/**: (Optional) Database schema definitions.
  - **migrates/**: Database migration files.

### Other

- **tests/**: Integration and unit tests.
- **logs/**: Application logs (default output directory).

## Getting Started

### Prerequisites

- Go 1.24 or higher
- Git
- (Optional) Docker for containerized database services

### Installation

1. Clone the repository:
   `+"```bash"+`
   git clone <repository-url>
   cd %s
   `+"```"+`

2. Install dependencies:
   `+"```bash"+`
   go mod tidy
   `+"```"+`

### Configuration

The application uses `+"`config.yaml`"+` for configuration.
Edit the configuration file to set your database connection, server port, etc.

`+"```bash"+`
vim config.yaml
`+"```"+`

### Running the Application

Start the server using:

`+"```bash"+`
go run cmd/main.go
`+"```"+`

The server will start on port **8080** (by default).
Access the health check at: `+"`http://localhost:8080/health`"+`

### Testing

Run all tests with:

`+"```bash"+`
go test ./...
`+"```"+`

## Development Guide

1. **Define Models**: Create your data models in `+"`data/model`"+` or `+"`data/schema`"+` (if using Ent).
2. **Repository**: Implement data access methods in `+"`data/repository`"+`.
3. **Service**: Implement business logic in `+"`service/`"+`, calling the repository.
4. **Handler**: Create HTTP handlers in `+"`handler/`"+` to map requests to service methods.
5. **Router**: Register your new handlers in `+"`internal/server/rest.go`"+` or `+"`internal/server/router.go`"+`.

## License

[Add License Information Here]
`, data.Name, data.Name)

	if err := utils.WriteTemplateFile(readmePath, readmeContent, nil); err != nil {
		fmt.Printf("Warning: failed to create README.md file: %v\n", err)
		// Just warn, don't stop the process
	}


	// Create Makefile for build support
	makefilePath := filepath.Join(basePath, "Makefile")
	makefileContent := fmt.Sprintf(`# Makefile for %s

# Binary name
BINARY_NAME=%s
OUTPUT_DIR=bin

# Version information
VERSION ?= $(shell git describe --tags --match "v*" --always 2>/dev/null || echo "v0.0.0")
BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
REVISION ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u '+%%Y-%%m-%%dT%%H:%%M:%%SZ')
GO_VERSION ?= $(shell go version | cut -d' ' -f3)

# Linker flags to set version information
LDFLAGS=-ldflags "\
	-X '%s/internal/version.Version=$(VERSION)' \
	-X '%s/internal/version.Branch=$(BRANCH)' \
	-X '%s/internal/version.Revision=$(REVISION)' \
	-X '%s/internal/version.BuiltAt=$(BUILD_TIME)' \
	-X '%s/internal/version.GoVersion=$(GO_VERSION)'"

.PHONY: all build run clean test help

all: build

## build: Build the application binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(OUTPUT_DIR)
	@go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME) ./cmd/$(BINARY_NAME)
	@echo "Build complete: $(OUTPUT_DIR)/$(BINARY_NAME)"

## run: Run the application
run:
	@go run $(LDFLAGS) ./cmd/$(BINARY_NAME)

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(OUTPUT_DIR)
	@go clean
	@echo "Clean complete"

## test: Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

## lint: Run linters
lint:
	@echo "Running linters..."
	@golangci-lint run || echo "golangci-lint not installed, skipping..."

## fmt: Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

## tidy: Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

## version: Show version information
version:
	@echo "Version:    $(VERSION)"
	@echo "Branch:     $(BRANCH)"
	@echo "Revision:   $(REVISION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Go Version: $(GO_VERSION)"

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@sed -n 's/^##//p' Makefile | column -t -s ':' | sed -e 's/^/ /'
`, data.Name, data.Name, data.PackagePath, data.PackagePath, data.PackagePath, data.PackagePath, data.PackagePath)

	if err := utils.WriteTemplateFile(makefilePath, makefileContent, nil); err != nil {
		fmt.Printf("Warning: failed to create Makefile: %v\n", err)
		// Just warn, don't stop the process
	}

	return nil

}
