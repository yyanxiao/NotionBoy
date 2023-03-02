package dao

import (
	"context"
	"notionboy/db/ent"
	"notionboy/db/ent/conversation"
	"notionboy/internal/pkg/db"

	"github.com/google/uuid"
)

func SaveConversation(ctx context.Context, userId uuid.UUID, instruction string) (*ent.Conversation, error) {
	return db.GetClient().Conversation.Create().
		SetInstruction(instruction).
		SetUserID(userId).
		SetUUID(uuid.New()).
		Save(ctx)
}

func GetConversation(ctx context.Context, conversationId uuid.UUID) (*ent.Conversation, error) {
	return db.GetClient().Conversation.Query().
		Where(conversation.UUIDEQ(conversationId)).
		Only(ctx)
}

func ListConversations(ctx context.Context, userId uuid.UUID, limit, offset int) ([]*ent.Conversation, error) {
	return db.GetClient().Conversation.Query().
		Where(conversation.UserIDEQ(userId)).
		Order(ent.Desc(conversation.FieldCreatedAt)).
		Limit(limit).
		Offset(offset).
		All(ctx)
}

func DeleteConversation(ctx context.Context, conversationId uuid.UUID) error {
	_, err := db.GetClient().Conversation.
		Delete().
		Where(conversation.UUIDEQ(conversationId)).
		Exec(ctx)
	return err
}
