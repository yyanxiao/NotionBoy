package schema

import (
	"notionboy/db/ent/schema/mixin"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// ChatHistory holds the schema definition for the ChatHistory entity.
// It is used to store the chat history of ChatGPT.
type ChatHistory struct {
	ent.Schema
}

// Fields of the Account.
func (ChatHistory) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_id").Comment("user id"),
		field.Int("conversation_idx").Comment("conversation idx for the whole chat history"),
		field.UUID("conversation_id", uuid.UUID{}).Comment("Conversation ID for the whole conversation"),
		field.String("message_id").Optional().Comment("Message ID inside the conversation"),
		field.Int("message_idx").Optional().Comment("Index of the message inside the conversation"),
		field.Text("request").Optional().Comment("Request of the conversation"),
		field.Text("response").Optional().Comment("Response of the conversation"),
		field.Int("token_usage").Optional().Comment("Token usage of the conversation"),
	}
}

func (ChatHistory) Edges() []ent.Edge {
	return nil
}

func (ChatHistory) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
		mixin.DeletedMixin{},
	}
}
