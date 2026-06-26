package generator

import (
	"fmt"

	"github.com/ncobase/cli/generator/templates"
)

func buildTemplateData(opts *Options) *templates.Data {
	return &templates.Data{
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
}

func buildRenderPlan(opts *Options, data *templates.Data) (*renderPlan, error) {
	render := newRenderPlan()

	if opts.Standalone {
		standalone, err := buildStandaloneRenderPlan(data)
		if err != nil {
			return nil, err
		}
		render.merge(standalone)
		return render, nil
	}

	mainTemplate := getMainTemplate(opts.Type)
	extension, err := buildExtensionRenderPlan(data, mainTemplate)
	if err != nil {
		return nil, err
	}
	render.merge(extension)

	if opts.WithCmd {
		cmdPlan, err := buildCmdRenderPlan(data)
		if err != nil {
			return nil, err
		}
		render.merge(cmdPlan)
	}

	return render, nil
}

func buildPlan(opts *Options, data *templates.Data, basePath string, render *renderPlan) *Plan {
	plan := &Plan{
		Name:         opts.Name,
		Type:         opts.Type,
		CustomDir:    opts.CustomDir,
		OutputPath:   opts.OutputPath,
		BasePath:     basePath,
		ModuleName:   opts.ModuleName,
		PackagePath:  data.PackagePath,
		Description:  getDesc(data),
		Standalone:   opts.Standalone,
		WithCmd:      opts.WithCmd || opts.Standalone,
		WithTest:     opts.WithTest,
		WithGRPC:     opts.WithGRPC,
		WithTracing:  opts.WithTracing,
		Database:     databasePlan(opts),
		Integrations: integrationPlan(opts),
		Directories:  render.directoryList(),
		Files:        render.fileList(),
	}

	if needsGoModule(opts) {
		plan.ModuleRequirements = moduleRequirementList(opts)
		plan.Operations = goModuleOperations(basePath, opts)
	}

	return plan
}

func databasePlan(opts *Options) DatabasePlan {
	orm := "none"
	switch {
	case opts.UseMongo:
		orm = "mongodb"
	case opts.UseEnt:
		orm = "ent"
	case opts.UseGorm:
		orm = "gorm"
	}

	driver := opts.DBDriver
	if driver == "" {
		driver = "none"
	}

	return DatabasePlan{
		ORM:    orm,
		Driver: driver,
	}
}

func integrationPlan(opts *Options) IntegrationPlan {
	plan := IntegrationPlan{}
	if opts.UseRedis {
		plan.Cache = append(plan.Cache, "redis")
	}
	if opts.UseElastic {
		plan.Search = append(plan.Search, "elasticsearch")
	}
	if opts.UseOpenSearch {
		plan.Search = append(plan.Search, "opensearch")
	}
	if opts.UseMeili {
		plan.Search = append(plan.Search, "meilisearch")
	}
	if opts.UseKafka {
		plan.Messaging = append(plan.Messaging, "kafka")
	}
	if opts.UseRabbitMQ {
		plan.Messaging = append(plan.Messaging, "rabbitmq")
	}
	if opts.UseS3Storage {
		plan.Storage = append(plan.Storage, "aws-s3")
	}
	if opts.UseMinio {
		plan.Storage = append(plan.Storage, "minio")
	}
	if opts.UseAliyun {
		plan.Storage = append(plan.Storage, "aliyun-oss")
	}
	if opts.WithGRPC {
		plan.Services = append(plan.Services, "grpc")
	}
	if opts.WithTracing {
		plan.Services = append(plan.Services, "opentelemetry")
	}
	return plan
}

func needsGoModule(opts *Options) bool {
	return opts.Standalone || opts.WithCmd
}

func successMessage(plan *Plan) string {
	if plan.Standalone {
		return fmt.Sprintf("Successfully generated standalone application %q in %s.", plan.Name, plan.Description)
	}
	return fmt.Sprintf("Successfully generated %q in %s.", plan.Name, plan.Description)
}
