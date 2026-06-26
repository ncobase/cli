package generator

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ncobase/cli/generator/templates"
	"github.com/ncobase/cli/utils"
	"github.com/ncobase/cli/version"
)

// initializeGoModule initializes a Go module for the generated code
func initializeGoModule(basePath string, data *templates.Data, opts *Options) error {
	goModPath := filepath.Join(basePath, "go.mod")

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

	var builder strings.Builder
	fmt.Fprintf(&builder, "module %s\n\ngo %s\n\nrequire (\n", data.PackagePath, version.DefaultGoVersion)
	for _, module := range sortedModules(requirements) {
		fmt.Fprintf(&builder, "\t%s %s\n", module, requirements[module])
	}
	builder.WriteString(")\n")

	if err := utils.WriteTemplateFile(goModPath, builder.String(), nil); err != nil {
		return fmt.Errorf("failed to create go.mod file: %w", err)
	}

	if opts.UseEnt {
		if err := runGoCommand(basePath, "go", "mod", "tidy", "-e"); err != nil {
			return err
		}
		if err := generateEntCode(basePath); err != nil {
			return err
		}
	}

	if err := runGoCommand(basePath, "go", "mod", "tidy"); err != nil {
		return err
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

func generateEntCode(basePath string) error {
	return runGoCommand(basePath,
		"go", "run", "-mod=mod", "entgo.io/ent/cmd/ent",
		"generate",
		"--feature", "sql/versioned-migration,sql/execquery,sql/upsert",
		"--target", "./data/ent",
		"./data/schema",
	)
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
