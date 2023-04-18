package schema

import (
	"notionboy/db/ent/schema/mixin"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Account holds the schema definition for the Account entity.
type WechatSession struct {
	ent.Schema
}

// Fields of the Account.
func (WechatSession) Fields() []ent.Field {
	return []ent.Field{
		field.Bytes("session").Comment("wechat login session"),
		field.String("dummy_user_id").Unique().Comment("dummy user_id"),
	}
}

func (WechatSession) Edges() []ent.Edge {
	return nil
}

func (WechatSession) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
		mixin.DeletedMixin{},
	}
}

func (WechatSession) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "wechat_session"},
	}
}
