package templates

func SchemaTemplate() string {
	return `package schema

import (
	"strings"

	"github.com/ncobase/ncore/data/databases/entgo/mixin"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
)

// Item holds the schema definition for the Item entity.
type Item struct {
	ent.Schema
}

// Annotations of the Item.
func (Item) Annotations() []schema.Annotation {
	table := strings.Join([]string{"ncse", "sys", "item"}, "_")
	return []schema.Annotation{
		entsql.Annotation{Table: table},
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
		entsql.WithComments(true),
	}
}

// Mixin of the Item.
func (Item) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.PrimaryKey,
		mixin.Name,
		mixin.Code,
		mixin.Status,
		mixin.Description,
		mixin.Operator,
		mixin.ExtraProps,
		mixin.CreatedAt,
		mixin.UpdatedAt,
	}
}

// Fields of the Item.
func (Item) Fields() []ent.Field {
	return nil
}

// Edges of the Item.
func (Item) Edges() []ent.Edge {
	return nil
}
`
}
