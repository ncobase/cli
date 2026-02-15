package generator

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed tmpl/**/*.tmpl
var templatesFS embed.FS

// Loader manages template loading and rendering
type Loader struct {
	templates *template.Template
	funcs     template.FuncMap
}

// NewLoader creates a new template loader
func NewLoader() (*Loader, error) {
	l := &Loader{
		funcs: defaultFuncMap(),
	}

	// Parse all templates from embedded FS
	tmpl := template.New("").Funcs(l.funcs)

	err := fs.WalkDir(templatesFS, "tmpl", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		content, err := templatesFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", path, err)
		}

		// Use relative path as template name
		name := strings.TrimPrefix(path, "tmpl/")
		name = strings.TrimSuffix(name, ".tmpl")

		_, err = tmpl.New(name).Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", name, err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	l.templates = tmpl
	return l, nil
}

// Render renders a template with the given data
func (l *Loader) Render(name string, data any) (string, error) {
	var buf strings.Builder

	tmpl := l.templates.Lookup(name)
	if tmpl == nil {
		return "", fmt.Errorf("template %s not found", name)
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to render template %s: %w", name, err)
	}

	return buf.String(), nil
}

// MustRender renders a template and panics on error
func (l *Loader) MustRender(name string, data any) string {
	result, err := l.Render(name, data)
	if err != nil {
		panic(err)
	}
	return result
}

// defaultFuncMap returns the default template functions
func defaultFuncMap() template.FuncMap {
	return template.FuncMap{
		// String functions
		"lower":     strings.ToLower,
		"upper":     strings.ToUpper,
		"trimSpace": strings.TrimSpace,
		"contains":  strings.Contains,
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,
		"join":      strings.Join,
		"split":     strings.Split,
		"replace":   strings.ReplaceAll,

		// Path functions
		"base": filepath.Base,
		"dir":  filepath.Dir,
		"ext":  filepath.Ext,

		// Logic functions
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
		"div": func(a, b int) int { return a / b },
		"mod": func(a, b int) int { return a % b },
		"not": func(b bool) bool { return !b },
		"and": func(a, b bool) bool { return a && b },
		"or":  func(a, b bool) bool { return a || b },

		// Formatting functions
		"sprintf": fmt.Sprintf,
		"quote":   func(s string) string { return fmt.Sprintf("%q", s) },

		// Database driver imports
		"driverImport": func(driver string) string {
			if driver == "" || driver == "none" {
				return ""
			}
			return fmt.Sprintf(`_ "github.com/ncobase/ncore/data/%s"`, driver)
		},

		// Conditional import generation
		"conditionalImport": func(condition bool, pkg string) string {
			if !condition {
				return ""
			}
			return fmt.Sprintf(`_ "%s"`, pkg)
		},
	}
}
