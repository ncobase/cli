package templates

import "fmt"

func GeneraterTemplate(name, extType, moduleName string) string {
	return fmt.Sprintf(`package %s

// Generate ent schema with versioned migrations
//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/versioned-migration --target ./data/ent {{ .PackagePath }}/data/schema

`, name)
}
