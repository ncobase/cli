package templates

import "fmt"

func GeneraterTemplate(name, extType, moduleName string) string {
	return fmt.Sprintf(`package %s

// Generate ent schema with versioned migrations and SQL helpers.
//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/versioned-migration,sql/execquery,sql/upsert --target ./data/ent ./data/schema

`, name)
}
