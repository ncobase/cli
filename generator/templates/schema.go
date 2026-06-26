package templates

func SchemaTemplate() string {
	return `package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Item holds the schema definition for the item entity.
type Item struct {
	ent.Schema
}

// Annotations returns database annotations for Item.
func (Item) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "items"},
		entsql.WithComments(true),
	}
}

// Fields returns item fields.
func (Item) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable(),
		field.String("name").
			NotEmpty().
			MaxLen(255),
		field.String("code").
			Default("").
			MaxLen(128),
		field.String("type").
			Default("").
			MaxLen(64),
		field.String("status").
			Default("active").
			MaxLen(64),
		field.String("description").
			Default(""),
		field.String("created_by").
			Default("").
			MaxLen(64),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.String("updated_by").
			Default("").
			MaxLen(64),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.Time("deleted_at").
			Optional().
			Nillable(),
		field.JSON("extras", map[string]any{}).
			Optional(),
	}
}

// Edges returns item edges.
func (Item) Edges() []ent.Edge {
	return nil
}
`
}
