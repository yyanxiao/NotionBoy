package schema

import (
	"notionboy/db/ent/schema/mixin"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Quota is the schema definition for the Quota entity.
// It is used to store the quota of the user.
type Quota struct {
	ent.Schema
}

// Fields of the Account.
func (Quota) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_id").Comment("user id"),
		field.Enum("category").Values("chatgpt"),
		field.Int("daily").Optional(),
		field.Int("monthly").Optional(),
		field.Int("yearly").Optional(),
		field.Int("daily_used").Optional(),
		field.Int("monthly_used").Optional(),
		field.Int("yearly_used").Optional(),
	}
}

func (Quota) Edges() []ent.Edge {
	return nil
}

func (Quota) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
		mixin.DeletedMixin{},
	}
}

func (Quota) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "category").Unique(),
	}
}
