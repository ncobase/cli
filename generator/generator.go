// Package generator provides code generation capabilities for ncobase CLI.
package generator

import (
	"fmt"
	"strings"

	"github.com/ncobase/cli/generator/templates"
	"github.com/ncobase/cli/utils"
)

// Options defines code generation options
type Options struct {
	Name        string
	Type        string
	ProjectType string
	CustomDir   string
	OutputPath  string
	ModuleName  string

	UseMongo bool
	UseEnt   bool
	UseGorm  bool

	WithCmd     bool
	WithTest    bool
	WithGRPC    bool
	WithTracing bool
	Group       string
	Standalone  bool
	DryRun      bool

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

// Generate builds and applies a generation plan. Dry-run mode returns the same plan without writing files.
func Generate(opts *Options) (*Result, error) {
	plan, render, err := prepareGeneration(opts)
	if err != nil {
		return nil, err
	}

	if opts.DryRun {
		return &Result{
			DryRun:  true,
			Applied: false,
			Message: fmt.Sprintf("Dry run complete for %q. No files were written.", plan.Name),
			Plan:    plan,
		}, nil
	}

	if len(plan.Conflicts) > 0 {
		return nil, fmt.Errorf("generation target has conflicts: %s", strings.Join(plan.Conflicts, "; "))
	}

	if err := writeRenderPlan(plan.BasePath, render); err != nil {
		return nil, err
	}

	if needsGoModule(opts) {
		if err := runGoModuleOperations(plan.BasePath, opts); err != nil {
			return nil, err
		}
	}

	return &Result{
		DryRun:  false,
		Applied: true,
		Message: successMessage(plan),
		Plan:    plan,
	}, nil
}

func prepareGeneration(opts *Options) (*Plan, *renderPlan, error) {
	if opts == nil {
		return nil, nil, fmt.Errorf("generation options are required")
	}

	opts.Name = strings.TrimSpace(opts.Name)
	opts.ProjectType = strings.TrimSpace(opts.ProjectType)
	opts.ModuleName = strings.TrimSpace(opts.ModuleName)
	opts.OutputPath = strings.TrimSpace(opts.OutputPath)
	opts.CustomDir = strings.TrimSpace(opts.CustomDir)
	opts.Group = strings.TrimSpace(opts.Group)

	if !utils.ValidateName(opts.Name) {
		return nil, nil, fmt.Errorf("invalid name %q: use letters, numbers, hyphens, or underscores, and start with a letter", opts.Name)
	}
	if opts.CustomDir != "" && !utils.ValidatePathSegment(opts.CustomDir) {
		return nil, nil, fmt.Errorf("invalid custom directory %q", opts.CustomDir)
	}
	if opts.Group != "" && !utils.ValidateName(opts.Group) {
		return nil, nil, fmt.Errorf("invalid group %q", opts.Group)
	}

	if err := normalizeOptions(opts); err != nil {
		return nil, nil, err
	}

	outputPath, err := resolveOutputPath(opts)
	if err != nil {
		return nil, nil, err
	}
	opts.OutputPath = outputPath
	opts.ModuleName = resolveModuleName(opts, outputPath)
	if strings.ContainsAny(opts.ModuleName, " \t\r\n") {
		return nil, nil, fmt.Errorf("module name %q must not contain whitespace", opts.ModuleName)
	}

	basePath := getBasePath(opts, outputPath)
	data := buildTemplateData(opts)

	render, err := buildRenderPlan(opts, data)
	if err != nil {
		return nil, nil, err
	}

	if needsGoModule(opts) {
		render.addFile("go.mod", buildGoModContent(data, opts))
	}

	plan := buildPlan(opts, data, basePath, render)
	if utils.PathExists(basePath) {
		plan.Conflicts = append(plan.Conflicts, fmt.Sprintf("directory %q already exists", basePath))
	}

	return plan, render, nil
}

func normalizeOptions(opts *Options) error {
	if opts.Standalone {
		projectType, err := normalizeProjectType(opts.ProjectType)
		if err != nil {
			return err
		}
		opts.ProjectType = projectType
	} else {
		opts.ProjectType = ""
	}

	opts.DBDriver = strings.ToLower(strings.TrimSpace(opts.DBDriver))

	if opts.DBDriver == "postgresql" {
		opts.DBDriver = "postgres"
	}
	if opts.DBDriver == "sqlite3" {
		opts.DBDriver = "sqlite"
	}
	if opts.DBDriver != "" && opts.DBDriver != "none" {
		switch opts.DBDriver {
		case "postgres", "mysql", "sqlite", "mongodb":
		default:
			return fmt.Errorf("unsupported database driver %q", opts.DBDriver)
		}
	}

	if opts.DBDriver == "mongodb" {
		opts.UseMongo = true
	}
	if opts.UseMongo && opts.DBDriver == "" {
		opts.DBDriver = "mongodb"
	}
	if opts.Standalone && opts.ProjectType == ProjectTypeModular && opts.DBDriver == "" && !opts.UseMongo {
		opts.DBDriver = "postgres"
	}
	if opts.Standalone && opts.ProjectType == ProjectTypeModular && (opts.UseEnt || opts.UseGorm) && opts.DBDriver == "" {
		opts.DBDriver = "postgres"
	}

	if opts.UseEnt && opts.UseGorm {
		return fmt.Errorf("use either --use-ent or --use-gorm, not both")
	}
	if opts.UseMongo && (opts.UseEnt || opts.UseGorm) {
		return fmt.Errorf("use --use-mongo without --use-ent or --use-gorm")
	}

	needsStandaloneData := (opts.Standalone && opts.ProjectType != ProjectTypeModular) || (!opts.Standalone && opts.WithCmd)
	if needsStandaloneData && opts.DBDriver != "" && opts.DBDriver != "none" && !opts.UseMongo && !opts.UseEnt && !opts.UseGorm {
		opts.UseEnt = true
	}
	if needsStandaloneData && !opts.UseMongo && !opts.UseEnt && !opts.UseGorm {
		opts.UseEnt = true
		opts.DBDriver = "sqlite"
	}
	if (opts.UseEnt || opts.UseGorm) && opts.DBDriver == "" {
		opts.DBDriver = "sqlite"
	}

	if opts.UseEnt || opts.UseGorm {
		switch opts.DBDriver {
		case "postgres", "mysql", "sqlite":
		default:
			return fmt.Errorf("database driver %q is not supported with SQL ORM; supported drivers are postgres, mysql, sqlite", opts.DBDriver)
		}
	}
	if opts.UseMongo && opts.DBDriver != "mongodb" {
		return fmt.Errorf("MongoDB projects must use --db mongodb")
	}
	if countEnabled(opts.UseS3Storage, opts.UseMinio, opts.UseAliyun) > 1 {
		return fmt.Errorf("choose only one storage driver: --use-s3, --use-minio, or --use-aliyun")
	}

	return nil
}

func getMainTemplate(typ string) func(string) string {
	switch typ {
	case "core":
		return templates.CoreTemplate
	case "biz":
		return templates.BusinessTemplate
	case "business":
		return templates.BusinessTemplate
	case "plugin":
		return templates.PluginTemplate
	default:
		return templates.BusinessTemplate
	}
}

func buildCmdRenderPlan(data *templates.Data) (*renderPlan, error) {
	plan := newRenderPlan()
	plan.addDir("cmd", "internal/server", "internal/middleware", "internal/version")

	files := map[string]string{
		"cmd/main.go":                             templates.CmdMainTemplate(data),
		"internal/server/server.go":               templates.ServerTemplate(data.PackagePath),
		"internal/server/http.go":                 templates.ServerHTTPTemplate(data.PackagePath),
		"internal/server/exts.go":                 templates.ServerExtsTemplate(data.PackagePath),
		"internal/middleware/cors.go":             templates.MiddlewareCORSTemplate(),
		"internal/middleware/security_headers.go": templates.MiddlewareSecurityHeadersTemplate(),
		"internal/middleware/trace.go":            templates.MiddlewareTraceTemplate(),
		"internal/middleware/logger.go":           templates.MiddlewareLoggerTemplate(),
		"internal/middleware/client_info.go":      templates.MiddlewareClientInfoTemplate(),
		"internal/version/version.go":             templates.VersionTemplate(),
	}

	for filePath, tmpl := range files {
		content, err := renderTemplateString(filePath, tmpl, data)
		if err != nil {
			return nil, fmt.Errorf("failed to render file %s: %v", filePath, err)
		}
		plan.addFile(filePath, content)
	}

	return plan, nil
}

func countEnabled(values ...bool) int {
	count := 0
	for _, value := range values {
		if value {
			count++
		}
	}
	return count
}
