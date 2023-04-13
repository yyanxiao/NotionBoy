package dao

import (
	"context"

	"notionboy/db/ent"
	"notionboy/db/ent/conversation"
	"notionboy/internal/pkg/db"

	"github.com/google/uuid"
)

func CreateConversation(ctx context.Context, id, userId uuid.UUID, instruction, title string) (*ent.Conversation, error) {
	uid := id
	if id == uuid.Nil {
		uid = uuid.New()
	}
	return db.GetClient().Conversation.Create().
		SetInstruction(instruction).
		SetTitle(title).
		SetUserID(userId).
		SetUUID(uid).
		Save(ctx)
}

func UpdateConversation(ctx context.Context, userId, id uuid.UUID, instruction, title string) (*ent.Conversation, error) {
	query := db.GetClient().Conversation.Create().
		SetUUID(id).
		SetUserID(userId)

	needUpdate := false

	if instruction != "" {
		query.SetInstruction(instruction)
		needUpdate = true
	}

	if title != "" {
		query.SetTitle(title)
		needUpdate = true
	}

	if needUpdate {
		if err := query.
			OnConflict().
			UpdateNewValues().
			Exec(ctx); err != nil {
			return nil, err
		}
	}

	return GetConversation(ctx, id)
}

func GetConversation(ctx context.Context, conversationId uuid.UUID) (*ent.Conversation, error) {
	return db.GetClient().Conversation.Query().
		Where(conversation.UUIDEQ(conversationId)).
		Only(ctx)
}

func ListConversations(ctx context.Context, userId uuid.UUID, limit, offset int) ([]*ent.Conversation, error) {
	if limit == 0 {
		limit = 10
	}
	return db.GetClient().Conversation.Query().
		Where(conversation.UserIDEQ(userId), conversation.TokenUsageGT(0)).
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

func IncrConversationUsedToken(cli *ent.Client, ctx context.Context, conversationId uuid.UUID, tokens int64) error {
	return cli.Conversation.
		Update().
		Where(conversation.UUIDEQ(conversationId)).
		AddTokenUsage(tokens).
		Exec(ctx)
}
