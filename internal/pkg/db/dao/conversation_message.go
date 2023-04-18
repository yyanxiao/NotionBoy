package dao

import (
	"context"

	"notionboy/db/ent"
	"notionboy/db/ent/conversationmessage"
	"notionboy/internal/pkg/db"

	"github.com/google/uuid"
)

func SaveConversationMessage(ctx context.Context, msg *ent.ConversationMessage) (*ent.ConversationMessage, error) {
	query := db.GetClient().ConversationMessage.Create().
		SetUserID(msg.UserID).
		SetConversationID(msg.ConversationID).
		SetRequest(msg.Request).
		SetResponse(msg.Response).
		SetTokenUsage(msg.TokenUsage).
		SetUUID(uuid.New())

	return query.Save(ctx)
}

func GetConversationMessage(ctx context.Context, ConversationMessageId uuid.UUID) (*ent.ConversationMessage, error) {
	return db.GetClient().ConversationMessage.Query().
		Where(conversationmessage.UUIDEQ(ConversationMessageId)).
		Only(ctx)
}

func ListConversationMessages(ctx context.Context, conversationId uuid.UUID, limit, offset int) ([]*ent.ConversationMessage, error) {
	if limit == 0 {
		limit = 10
	}
	return db.GetClient().ConversationMessage.Query().
		Where(conversationmessage.ConversationID(conversationId)).
		Order(ent.Asc(conversationmessage.FieldCreatedAt)).
		Limit(limit).
		Offset(offset).
		All(ctx)
}

func DeleteConversationMessage(ctx context.Context, ConversationMessageId uuid.UUID) error {
	_, err := db.GetClient().ConversationMessage.
		Delete().
		Where(conversationmessage.UUIDEQ(ConversationMessageId)).
		Exec(ctx)
	return err
}
