package generator

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ncobase/cli/generator/templates"
	"github.com/ncobase/cli/utils"
	"github.com/ncobase/cli/version"
)

// initializeGoModule initializes a Go module for the generated code
func initializeGoModule(basePath string, data *templates.Data, opts *Options) error {
	goModPath := filepath.Join(basePath, "go.mod")

	var builder strings.Builder
	fmt.Fprintf(&builder, "module %s\n\ngo %s\n\nrequire (\n\tgithub.com/gin-gonic/gin %s\n\tgithub.com/google/uuid %s\n)\n",
		data.PackagePath, version.DefaultGoVersion, version.GinVersion, version.UUIDVersion)

	ncoreDeps := version.GetNcoreDeps()

	if opts.UseMongo {
		fmt.Fprintf(&builder, "\nrequire go.mongodb.org/mongo-driver %s\n", version.MongoVersion)
	}
	if opts.UseEnt {
		fmt.Fprintf(&builder, "\nrequire entgo.io/ent %s\n", version.EntVersion)
	}
	if opts.UseGorm {
		fmt.Fprintf(&builder, "\nrequire gorm.io/gorm %s\n", version.GormVersion)
	}

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

	builder.WriteString("\n// Replace directives for local ncore development\n")
	builder.WriteString("// Remove these lines and run 'go mod tidy' when ncore packages are published\n")
	for _, dep := range ncoreDeps {
		parts := strings.Split(dep, "/")
		moduleName := parts[len(parts)-1]
		fmt.Fprintf(&builder, "// replace %s => ../../ncore/%s\n", dep, moduleName)
	}

	if err := utils.WriteTemplateFile(goModPath, builder.String(), nil); err != nil {
		return fmt.Errorf("failed to create go.mod file: %v", err)
	}

	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = basePath
	if err := tidyCmd.Run(); err != nil {
		fmt.Printf("Warning: failed to run 'go mod tidy': %v\n", err)
	}

	if opts.UseEnt {
		schemaDir := filepath.Join(basePath, "data/schema")
		if err := utils.EnsureDir(schemaDir); err != nil {
			fmt.Printf("Warning: failed to create ent schema directory: %v\n", err)
		}
	}

	return nil
}
