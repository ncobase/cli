package generator

import (
	"bytes"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/ncobase/cli/generator/templates"
	"github.com/ncobase/cli/version"
)

func buildModuleRequirements(opts *Options) map[string]string {
	requirements := map[string]string{
		"github.com/gin-gonic/gin":            version.GinVersion,
		"github.com/google/uuid":              version.UUIDVersion,
		"github.com/prometheus/client_golang": version.PrometheusVersion,
		"github.com/sirupsen/logrus":          version.LogrusVersion,
	}

	if opts.UseMongo {
		requirements["go.mongodb.org/mongo-driver/v2"] = version.MongoVersion
	}
	if opts.UseEnt {
		requirements["entgo.io/ent"] = version.EntVersion
	}
	if opts.UseGorm {
		requirements["gorm.io/gorm"] = version.GormVersion
		addGormDriverRequirements(requirements, opts.DBDriver)
	}
	if opts.WithTest {
		requirements["github.com/stretchr/testify"] = "v1.11.1"
	}

	for _, dep := range version.GetNcoreDeps() {
		requirements[dep] = version.NcoreVersion(dep)
	}
	if opts.DBDriver != "" && opts.DBDriver != "none" {
		requirements["github.com/ncobase/ncore/data/"+opts.DBDriver] = version.NcoreVersion("github.com/ncobase/ncore/data/" + opts.DBDriver)
	}
	if opts.UseRedis {
		requirements["github.com/ncobase/ncore/data/redis"] = version.NcoreVersion("github.com/ncobase/ncore/data/redis")
	}
	if opts.UseElastic {
		requirements["github.com/ncobase/ncore/data/elasticsearch"] = version.NcoreVersion("github.com/ncobase/ncore/data/elasticsearch")
	}
	if opts.UseOpenSearch {
		requirements["github.com/ncobase/ncore/data/opensearch"] = version.NcoreVersion("github.com/ncobase/ncore/data/opensearch")
	}
	if opts.UseMeili {
		requirements["github.com/ncobase/ncore/data/meilisearch"] = version.NcoreVersion("github.com/ncobase/ncore/data/meilisearch")
	}
	if opts.UseKafka {
		requirements["github.com/ncobase/ncore/data/kafka"] = version.NcoreVersion("github.com/ncobase/ncore/data/kafka")
	}
	if opts.UseRabbitMQ {
		requirements["github.com/ncobase/ncore/data/rabbitmq"] = version.NcoreVersion("github.com/ncobase/ncore/data/rabbitmq")
	}
	if opts.UseS3Storage || opts.UseMinio || opts.UseAliyun {
		requirements["github.com/ncobase/ncore/oss"] = version.NcoreVersion("github.com/ncobase/ncore/oss")
	}

	return requirements
}

func buildGoModContent(data *templates.Data, opts *Options) string {
	requirements := buildModuleRequirements(opts)

	var builder strings.Builder
	fmt.Fprintf(&builder, "module %s\n\ngo %s\n\nrequire (\n", data.PackagePath, version.DefaultGoVersion)
	for _, module := range sortedModules(requirements) {
		fmt.Fprintf(&builder, "\t%s %s\n", module, requirements[module])
	}
	builder.WriteString(")\n")

	return builder.String()
}

func moduleRequirementList(opts *Options) []ModuleRequirement {
	requirements := buildModuleRequirements(opts)
	modules := sortedModules(requirements)
	result := make([]ModuleRequirement, 0, len(modules))
	for _, module := range modules {
		result = append(result, ModuleRequirement{
			Module:  module,
			Version: requirements[module],
		})
	}
	return result
}

func goModuleOperations(basePath string, opts *Options) []Operation {
	operations := make([]Operation, 0, 3)
	if opts.UseEnt {
		operations = append(operations,
			Operation{
				Name:       "go mod tidy for Ent code generation",
				Command:    "go",
				Args:       []string{"mod", "tidy", "-e"},
				WorkingDir: basePath,
				Outputs:    []string{"go.sum"},
			},
			Operation{
				Name:       "generate Ent client",
				Command:    "go",
				Args:       []string{"run", "-mod=mod", "entgo.io/ent/cmd/ent", "generate", "--feature", "sql/versioned-migration,sql/execquery,sql/upsert", "--target", "./data/ent", "./data/schema"},
				WorkingDir: basePath,
				Outputs:    []string{"data/ent/*", "go.sum"},
			},
		)
	}

	operations = append(operations, Operation{
		Name:       "go mod tidy",
		Command:    "go",
		Args:       []string{"mod", "tidy"},
		WorkingDir: basePath,
		Outputs:    []string{"go.sum"},
	})

	return operations
}

func runGoModuleOperations(basePath string, opts *Options) error {
	for _, operation := range goModuleOperations(basePath, opts) {
		if err := runGoCommand(basePath, operation.Command, operation.Args...); err != nil {
			return err
		}
	}

	return nil
}

func addGormDriverRequirements(requirements map[string]string, driver string) {
	switch driver {
	case "postgres", "postgresql":
		requirements["gorm.io/driver/postgres"] = "v1.6.0"
	case "mysql":
		requirements["gorm.io/driver/mysql"] = "v1.6.0"
	case "sqlite", "sqlite3":
		requirements["gorm.io/driver/sqlite"] = "v1.6.0"
	case "sqlserver", "mssql":
		requirements["gorm.io/driver/sqlserver"] = "v1.6.3"
	}
}

func sortedModules(requirements map[string]string) []string {
	modules := make([]string, 0, len(requirements))
	for module := range requirements {
		modules = append(modules, module)
	}
	sort.Strings(modules)
	return modules
}

func runGoCommand(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s %s failed: %w\n%s", name, strings.Join(args, " "), err, strings.TrimSpace(output.String()))
	}
	return nil
}
