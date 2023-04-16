package schema

import (
	"notionboy/db/ent/schema/mixin"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Prompt struct {
	ent.Schema
}

func (Prompt) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("uuid", uuid.UUID{}).Unique().Immutable().Comment("UUID"),
		field.UUID("user_id", uuid.UUID{}).Comment("user id"),
		field.String("act").Comment("role name"),
		field.String("prompt").Comment("prompt text"),
		field.Bool("is_custom").Comment("is user custom prompt"),
	}
}

func (Prompt) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
		mixin.DeletedMixin{},
	}
}
