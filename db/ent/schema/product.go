package schema

import (
	"notionboy/db/ent/schema/mixin"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Product struct {
	ent.Schema
}

func (Product) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("uuid", uuid.UUID{}).Unique().Immutable().Comment("UUID"),
		field.String("name").Default("Free").Comment("product name"),
		field.Text("description").Comment("product description"),
		field.Float("price").Comment("product price"),
		field.Int64("token").Default(10000).Comment("Contained number of OpenAI token"),
		field.Int64("storage").Default(100).Comment("Contained size of S3 storage, unit: MB"),
	}
}

func (Product) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
		mixin.DeletedMixin{},
	}
}
