package schema

import (
	"notionboy/db/ent/schema/mixin"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// ChatHistory holds the schema definition for the ChatHistory entity.
// It is used to store the chat history of ChatGPT.
type Conversation struct {
	ent.Schema
}

// Fields of the Account.
func (Conversation) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("uuid", uuid.UUID{}).Unique().Immutable().Comment("UUID"),
		field.UUID("user_id", uuid.UUID{}).Comment("user id"),
		field.Text("instruction").Optional().Comment("Instructions for the conversation"),
	}
}

func (Conversation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("conversation_messages", ConversationMessage.Type),
	}
}

func (Conversation) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
		mixin.DeletedMixin{},
	}
}

func (Conversation) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "created_at"), // index for user_id and created_at for pagination
	}
}
