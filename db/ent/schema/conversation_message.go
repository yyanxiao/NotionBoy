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
type ConversationMessage struct {
	ent.Schema
}

// Fields of the Account.
func (ConversationMessage) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("uuid", uuid.UUID{}).Unique().Immutable().Comment("UUID"),
		field.UUID("user_id", uuid.UUID{}).Comment("user id"),
		field.UUID("conversation_id", uuid.UUID{}).Comment("Conversation ID for the conversation"),
		field.Text("request").Optional().Comment("Request of the message"),
		field.Text("response").Optional().Comment("Response of the message"),
		field.Int64("token_usage").Optional().Comment("Token usage of the message in the conversation"),
		field.String("model").Default("gpt-3.5-turbo").Comment("Model used for the message"),
	}
}

func (ConversationMessage) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("conversations", Conversation.Type).Ref("conversation_messages").Unique(),
	}
}

func (ConversationMessage) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
		mixin.DeletedMixin{},
	}
}

func (ConversationMessage) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "conversation_id", "created_at"), // for list messages in a conversation
	}
}
