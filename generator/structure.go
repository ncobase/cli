package generator

import (
	"fmt"

	"github.com/ncobase/cli/generator/templates"
)

func buildExtensionRenderPlan(data *templates.Data, mainTemplate func(string) string) (*renderPlan, error) {
	plan := newRenderPlan()
	directories := []string{
		"data", "data/repository",
		"handler", "router", "service", "structs",
	}
	if data.UseEnt {
		directories = append(directories, "data/schema")
	}
	if data.UseGorm {
		directories = append(directories, "data/model")
	}
	if data.WithTest {
		directories = append(directories, "tests")
	}
	plan.addDir(directories...)

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
		"handler/provider.go":           templates.HandlerTemplate(data.Name, data.ExtType, data.ModuleName),
		"router/provider.go":            templates.RouterTemplate(),
		"service/provider.go":           templates.ServiceTemplate(data.Name, data.ExtType, data.ModuleName),
		"structs/structs.go":            templates.StructsTemplate(),
	}

	if data.UseEnt {
		files["data/schema/item.go"] = templates.SchemaTemplate()
		files["generate.go"] = templates.GeneraterTemplate(data.Name, data.ExtType, data.ModuleName)
	}
	if data.UseGorm {
		files["data/model/item.go"] = templates.GormItemModelTemplate()
	}

	if data.WithTest {
		files["tests/ext_test.go"] = templates.ExtTestTemplate(data.Name, data.ExtType, data.ModuleName)
		files["tests/handler_test.go"] = templates.HandlerTestTemplate(data.Name, data.ExtType, data.ModuleName)
		files["tests/service_test.go"] = templates.ServiceTestTemplate(data.Name, data.ExtType, data.ModuleName)
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
