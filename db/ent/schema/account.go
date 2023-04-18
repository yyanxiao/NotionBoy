package schema

import (
	"notionboy/db/ent/schema/mixin"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Account holds the schema definition for the Account entity.
type Account struct {
	ent.Schema
}

// Fields of the Account.
func (Account) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("uuid", uuid.UUID{}).Unique().Optional().Comment("UUID"),
		field.String("user_id").Comment("user id").Comment("For wechat, it's openid, for telegram, it's telegram user id, for oauth, it's oauth user email"),
		field.Enum("user_type").Values("wechat", "telegram", "github", "google", "twitter", "microsoft").Optional().Default("wechat"),
		field.String("database_id").Optional().Comment("Notion Database ID"),
		field.String("access_token").Optional().Sensitive().Comment("Notion Access Token"),
		field.String("notion_user_id").Optional().Comment("Notion User ID"),
		field.String("notion_user_email").Optional().Comment("Notion User Email"),
		field.Bool("is_latest_schema").Default(false).Comment("If not the latest schema, need update notion page properies"),
		field.Bool("is_openai_api_user").Default(false).Comment("Dose this user can use openai API instead of reverse session"),
		field.String("openai_api_key").Optional().Sensitive().Comment("OpenAI API Key"),
		field.UUID("api_key", uuid.UUID{}).Optional().Comment("API Key"),
		field.Bool("is_admin").Default(false).Comment("Is admin user"),
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

func (Account) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "user_type").Unique(),
	}
}
