package dao

import (
	"context"
	"notionboy/db/ent"
	"notionboy/db/ent/chathistory"
	"notionboy/internal/pkg/db"

	"github.com/google/uuid"
)

// QueryChatHistory Get ChatHistory by user id and order by created_at
// 1. select latest conversation history by limit
// 2. get all messages from the conversations in step 1
// 3. return all messages order by conversation_idx and message_idx
func QueryChatHistory(ctx context.Context, userID int, limit int) ([]*ent.ChatHistory, error) {
	histories, err := db.GetClient().ChatHistory.Query().
		Select(chathistory.FieldConversationID).
		Where(chathistory.UserIDEQ(userID), chathistory.ConversationIdxEQ(0)).
		Order(ent.Desc(chathistory.FieldCreatedAt)).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, err
	}
	conversationIds := make([]uuid.UUID, 0, len(histories))
	for _, history := range histories {
		conversationIds = append(conversationIds, history.ConversationID)
	}
	return db.GetClient().ChatHistory.Query().
		Where(chathistory.UserIDEQ(userID), chathistory.ConversationIDIn(conversationIds...)).
		Order(ent.Desc(chathistory.FieldConversationIdx), ent.Desc(chathistory.FieldMessageIdx)).
		All(ctx)
}

// SaveChatHistory Save ChatHistory
func SaveChatHistory(ctx context.Context, h *ent.ChatHistory) error {
	query := db.GetClient().ChatHistory.Create().
		SetUserID(h.UserID).
		SetMessageID(h.MessageID).
		SetConversationID(h.ConversationID).
		SetRequest(h.Request).
		SetResponse(h.Response)

	if h.ConversationIdx == 0 {
		maxIdx, err := QueryMaxConversationIdx(ctx, h.UserID)
		if err != nil {
			return err
		}
		h.ConversationIdx = maxIdx + 1
		query.SetConversationIdx(h.ConversationIdx)
	} else {
		query.SetConversationIdx(h.ConversationIdx)
	}

	if h.MessageIdx == 0 {
		maxMessageIdx, err := QueryMaxMessageIdx(ctx, h.UserID, h.ConversationID)
		if err != nil {
			return err
		}
		h.MessageIdx = maxMessageIdx + 1
		query.SetMessageIdx(h.MessageIdx)
	} else {
		query.SetMessageIdx(h.MessageIdx)
	}

	return query.
		OnConflict().
		UpdateNewValues().
		Exec(ctx)
}

// QueryMaxConversationIdx Get Max HistoryIdx by user id
func QueryMaxConversationIdx(ctx context.Context, userID int) (int, error) {
	history, err := db.GetClient().ChatHistory.Query().
		Select(chathistory.FieldConversationIdx).
		Where(chathistory.UserIDEQ(userID)).
		Order(ent.Desc(chathistory.FieldConversationIdx)).
		Limit(1).
		Only(ctx)
	if err != nil {
		return 0, err
	}
	return history.ConversationIdx, nil
}

// QueryMaxMessageIdx Get Max message idx by user id and conversation id
func QueryMaxMessageIdx(ctx context.Context, userID int, conversationID uuid.UUID) (int, error) {
	history, err := db.GetClient().ChatHistory.Query().
		Select(chathistory.FieldMessageIdx).
		Where(chathistory.UserIDEQ(userID), chathistory.ConversationIDEQ(conversationID)).
		Order(ent.Desc(chathistory.FieldMessageIdx)).
		Limit(1).
		Only(ctx)
	if err != nil {
		return 0, err
	}
	return history.MessageIdx, nil
}
