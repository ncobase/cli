// Package generator provides code generation capabilities for ncobase CLI.
package generator

import (
	"fmt"
	"path/filepath"

	"github.com/ncobase/cli/generator/templates"
	"github.com/ncobase/cli/utils"
)

// Options defines code generation options
type Options struct {
	Name       string
	Type       string
	CustomDir  string
	OutputPath string
	ModuleName string

	UseMongo bool
	UseEnt   bool
	UseGorm  bool

	WithCmd     bool
	WithTest    bool
	WithGRPC    bool
	WithTracing bool
	Group       string
	Standalone  bool

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

// DefaultOptions returns default options
func DefaultOptions() *Options {
	return &Options{
		Type: "custom",
	}
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

	outputPath, err := resolveOutputPath(opts)
	if err != nil {
		return err
	}
	opts.OutputPath = outputPath
	opts.ModuleName = resolveModuleName(opts, outputPath)

	basePath := getBasePath(opts, outputPath)
	if utils.PathExists(basePath) {
		return fmt.Errorf("directory '%s' already exists", basePath)
	}

	if err := utils.EnsureDir(basePath); err != nil {
		return fmt.Errorf("failed to create base directory: %v", err)
	}

	data := &templates.Data{
		Name:          opts.Name,
		Type:          opts.Type,
		UseMongo:      opts.UseMongo,
		UseEnt:        opts.UseEnt,
		UseGorm:       opts.UseGorm,
		WithTest:      opts.WithTest,
		WithCmd:       opts.WithCmd || opts.Standalone,
		WithGRPC:      opts.WithGRPC,
		WithTracing:   opts.WithTracing,
		Standalone:    opts.Standalone,
		Group:         opts.Group,
		ExtType:       opts.Type,
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

	if opts.Standalone {
		if err := createStandaloneStructure(basePath, data); err != nil {
			return err
		}
		if err := initializeGoModule(basePath, data, opts); err != nil {
			fmt.Printf("Warning: failed to initialize Go module: %v\n", err)
		}
		fmt.Printf("Successfully generated standalone application '%s' in %s\n", data.Name, getDesc(data))
		return nil
	}

	mainTemplate := getMainTemplate(opts.Type)
	if err := createStructure(basePath, data, mainTemplate); err != nil {
		return err
	}

	if opts.WithCmd {
		if err := createCmdStructure(basePath, data); err != nil {
			return err
		}
		if err := initializeGoModule(basePath, data, opts); err != nil {
			fmt.Printf("Warning: failed to initialize Go module: %v\n", err)
		}
	}

	fmt.Printf("Successfully generated '%s' in %s\n", data.Name, getDesc(data))
	return nil
}

func getMainTemplate(typ string) func(string) string {
	switch typ {
	case "core":
		return templates.CoreTemplate
	case "business":
		return templates.BusinessTemplate
	case "plugin":
		return templates.PluginTemplate
	default:
		return templates.BusinessTemplate
	}
}

func createCmdStructure(basePath string, data *templates.Data) error {
	dirs := []string{"cmd", "internal/server", "internal/middleware", "internal/version"}
	for _, dir := range dirs {
		if err := utils.EnsureDir(filepath.Join(basePath, dir)); err != nil {
			return fmt.Errorf("failed to create %s: %v", dir, err)
		}
	}

	files := map[string]string{
		"cmd/main.go":                         templates.CmdMainTemplate(data),
		"internal/server/server.go":           templates.ServerTemplate(data.PackagePath),
		"internal/server/http.go":             templates.ServerHTTPTemplate(data.PackagePath),
		"internal/server/exts.go":             templates.ServerExtsTemplate(data.PackagePath),
		"internal/middleware/cors.go":         templates.MiddlewareCORSTemplate(),
		"internal/middleware/security_headers.go": templates.MiddlewareSecurityHeadersTemplate(),
		"internal/middleware/trace.go":        templates.MiddlewareTraceTemplate(),
		"internal/middleware/logger.go":       templates.MiddlewareLoggerTemplate(),
		"internal/middleware/client_info.go":  templates.MiddlewareClientInfoTemplate(),
		"internal/version/version.go":         templates.VersionTemplate(),
	}

	for filePath, tmpl := range files {
		if err := utils.WriteTemplateFile(filepath.Join(basePath, filePath), tmpl, data); err != nil {
			return fmt.Errorf("failed to create file %s: %v", filePath, err)
		}
	}

	return nil
}
