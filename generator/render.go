package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"text/template"
)

type renderPlan struct {
	directories []string
	files       map[string]string
}

func newRenderPlan() *renderPlan {
	return &renderPlan{
		files: make(map[string]string),
	}
}

func (p *renderPlan) addDir(paths ...string) {
	p.directories = append(p.directories, paths...)
}

func (p *renderPlan) addFile(path, content string) {
	p.files[path] = content
}

func (p *renderPlan) merge(other *renderPlan) {
	if other == nil {
		return
	}
	p.directories = append(p.directories, other.directories...)
	for path, content := range other.files {
		p.files[path] = content
	}
}

func (p *renderPlan) directoryList() []string {
	return sortedUnique(p.directories)
}

func (p *renderPlan) fileList() []string {
	files := make([]string, 0, len(p.files))
	for path := range p.files {
		files = append(files, path)
	}
	sort.Strings(files)
	return files
}

func writeRenderPlan(basePath string, plan *renderPlan) error {
	for _, dir := range plan.directoryList() {
		if err := os.MkdirAll(filepath.Join(basePath, dir), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	for _, filePath := range plan.fileList() {
		fullPath := filepath.Join(basePath, filePath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %v", filePath, err)
		}
		if err := os.WriteFile(fullPath, []byte(plan.files[filePath]), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %v", filePath, err)
		}
	}

	return nil
}

func renderTemplateString(name, content string, data any) (string, error) {
	tmpl, err := template.New(name).Parse(content)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func sortedUnique(values []string) []string {
	if len(values) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}
