package conversation

import (
	"context"
	"strings"

	"notionboy/db/ent"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/logger"

	"github.com/google/uuid"
	"google.golang.org/grpc/status"
)

func (m *conversationMgr) CreateConversation(ctx context.Context, acc *ent.Account, instruction string) (*ConversationDTO, error) {
	if instruction == "" {
		instruction = DEFAULT_INSTRUCTION
	}
	res, err := dao.SaveConversation(ctx, acc.UUID, instruction)
	return ConversationDTOFromDB(res), err
}

func (m *conversationMgr) GetConversation(ctx context.Context, acc *ent.Account, Id string) (*ConversationDTO, error) {
	id, err := uuid.Parse(Id)
	if err != nil {
		return nil, status.Errorf(400, "invalid id, %s", err.Error())
	}
	res, err := dao.GetConversation(ctx, id)
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
	id, err := uuid.Parse(Id)
	if err != nil {
		return status.Errorf(400, "invalid id %s", err.Error())
	}
	return dao.DeleteConversation(ctx, id)
}

func (m *conversationMgr) CreateConversationMessage(ctx context.Context, acc *ent.Account, conversationId, request string) (*ConversationMessageDTO, error) {
	id, err := uuid.Parse(conversationId)
	if err != nil {
		return nil, status.Errorf(400, "invalid id %s", err.Error())
	}
	conversation, err := dao.GetConversation(ctx, id)
	logger.SugaredLogger.Debugw("get conversation", "conversation", conversation, "err", err)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			logger.SugaredLogger.Debugw("conversation not found, create a new one", "conversationId", conversationId)
			conversationDTO, err := m.CreateConversation(ctx, acc, "")
			if err != nil {
				return nil, err
			}
			logger.SugaredLogger.Debugw("create conversation", "conversationDTO", conversationDTO)
			conversation = conversationDTO.ToDB()
		} else {
			return nil, err
		}
	}

	apiClient := NewApiClient(acc.OpenaiAPIKey)
	message, err := apiClient.ChatWithHistory(ctx, acc, conversation.Instruction, conversationId, request)
	logger.SugaredLogger.Debugw("chat with history", "message", message, "err", err)
	if err != nil {
		logger.SugaredLogger.Debugw("chat with history error", "err", err)
		return nil, err
	}
	return ConversationMessageDTOFromDB(message), nil
}

func (m *conversationMgr) GetConversationMessage(ctx context.Context, acc *ent.Account, conversationId, messageId string) (*ConversationMessageDTO, error) {
	id, err := uuid.Parse(messageId)
	if err != nil {
		return nil, status.Errorf(400, "invalid id %s", err.Error())
	}
	res, err := dao.GetConversationMessage(ctx, id)
	return ConversationMessageDTOFromDB(res), err
}

func (m *conversationMgr) ListConversationMessages(ctx context.Context, acc *ent.Account, conversationId string, limit, offset int) ([]*ConversationMessageDTO, error) {
	id, err := uuid.Parse(conversationId)
	if err != nil {
		return nil, status.Errorf(400, "invalid id %s", err.Error())
	}
	res, err := dao.ListConversationMessages(ctx, id, limit, offset)
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
	id, err := uuid.Parse(messageId)
	if err != nil {
		return status.Errorf(400, "invalid id %s", err.Error())
	}
	return dao.DeleteConversationMessage(ctx, id)
}
