package generator

import (
	"fmt"
	"path/filepath"

	"github.com/ncobase/cli/generator/templates"
	"github.com/ncobase/cli/utils"
)

// createStructure creates the extension structure
func createStructure(basePath string, data *templates.Data, mainTemplate func(string) string) error {
	if err := utils.EnsureDir(basePath); err != nil {
		return fmt.Errorf("failed to create base directory: %v", err)
	}

	directories := []string{
		"data", "data/repository", "data/schema",
		"handler", "service", "structs",
	}
	if data.WithTest {
		directories = append(directories, "tests")
	}

	for _, dir := range directories {
		if err := utils.EnsureDir(filepath.Join(basePath, dir)); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

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

	if data.UseEnt {
		files["generate.go"] = templates.GeneraterTemplate(data.Name, data.ExtType, data.ModuleName)
	}

	if data.WithTest {
		files["tests/ext_test.go"] = templates.ExtTestTemplate(data.Name, data.ExtType, data.ModuleName)
		files["tests/handler_test.go"] = templates.HandlerTestTemplate(data.Name, data.ExtType, data.ModuleName)
		files["tests/service_test.go"] = templates.ServiceTestTemplate(data.Name, data.ExtType, data.ModuleName)
	}

	for filePath, tmpl := range files {
		if err := utils.WriteTemplateFile(filepath.Join(basePath, filePath), tmpl, data); err != nil {
			return fmt.Errorf("failed to create file %s: %v", filePath, err)
		}
	}

	return nil
}
