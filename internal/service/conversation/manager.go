package conversation

import (
	"context"
	"notionboy/db/ent"
	"notionboy/internal/pkg/db/dao"

	"github.com/google/uuid"
)

func (m *conversationMgr) CreateConversation(ctx context.Context, acc *ent.Account, instruction string) (*ConversationDTO, error) {
	if instruction == "" {
		instruction = DEFAULT_INSTRUCTION
	}
	res, err := dao.SaveConversation(ctx, acc.UUID, instruction)
	return ConversationDTOFromDB(res), err
}

func (m *conversationMgr) GetConversation(ctx context.Context, acc *ent.Account, Id string) (*ConversationDTO, error) {
	res, err := dao.GetConversation(ctx, uuid.MustParse(Id))
	return ConversationDTOFromDB(res), err
}

func (m *conversationMgr) ListConversations(ctx context.Context, acc *ent.Account, limit, offset int) ([]*ConversationDTO, error) {
	res, err := dao.ListConversations(ctx, acc.UUID, limit, offset)
	if err != nil {
		return nil, err
	}
	ret := make([]*ConversationDTO, 0, len(res))
	for _, c := range res {
		ret = append(ret, ConversationDTOFromDB(c))
	}
	return ret, nil
}

func (m *conversationMgr) DeleteConversation(ctx context.Context, acc *ent.Account, Id string) error {
	return dao.DeleteConversation(ctx, uuid.MustParse(Id))
}

func (m *conversationMgr) CreateConversationMessage(ctx context.Context, acc *ent.Account, conversationId, request string) (*ConversationMessageDTO, error) {
	conversation, err := dao.GetConversation(ctx, uuid.MustParse(conversationId))
	if err != nil {
		return nil, err
	}

	apiClient := NewApiClient(acc.OpenaiAPIKey)
	message, err := apiClient.ChatWithHistory(ctx, acc, conversation.Instruction, conversationId, request)
	if err != nil {
		return nil, err
	}
	return ConversationMessageDTOFromDB(message), nil
}

func (m *conversationMgr) GetConversationMessage(ctx context.Context, acc *ent.Account, conversationId, messageId string) (*ConversationMessageDTO, error) {
	res, err := dao.GetConversationMessage(ctx, uuid.MustParse(messageId))
	return ConversationMessageDTOFromDB(res), err
}

func (m *conversationMgr) ListConversationMessages(ctx context.Context, acc *ent.Account, conversationId string, limit, offset int) ([]*ConversationMessageDTO, error) {
	res, err := dao.ListConversationMessages(ctx, uuid.MustParse(conversationId), limit, offset)
	if err != nil {
		return nil, err
	}
	ret := make([]*ConversationMessageDTO, 0, len(res))
	for _, c := range res {
		ret = append(ret, ConversationMessageDTOFromDB(c))
	}
	return ret, nil
}

func (m *conversationMgr) DeleteConversationMessage(ctx context.Context, acc *ent.Account, conversationId, messageId string) error {
	return dao.DeleteConversationMessage(ctx, uuid.MustParse(messageId))
}
