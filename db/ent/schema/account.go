package schema

import (
	"notionboy/db/ent/schema/mixin"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Account holds the schema definition for the Account entity.
type Account struct {
	ent.Schema
}

// Fields of the Account.
func (Account) Fields() []ent.Field {
	return []ent.Field{
		field.String("user_id").Unique().Comment("wechat user id"),
		field.Enum("user_type").Values("wechat", "telegram").Optional().Default("wechat"),
		field.String("database_id").Unique().NotEmpty().Comment("Notion Database ID"),
		field.String("access_token").Sensitive().NotEmpty().Comment("Notion Access Token"),
		field.String("notion_user_id").Optional().Comment("Notion User ID"),
		field.String("notion_user_email").Optional().Comment("Notion User Email"),
		field.Bool("is_latest_schema").Default(false).Comment("If not the latest schema, need update notion page properies"),
	}
}

func (Account) Edges() []ent.Edge {
	return nil
}

func (Account) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
		mixin.DeletedMixin{},
	}
}
