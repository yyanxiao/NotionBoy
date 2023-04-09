package conversation

import (
	"context"

	"notionboy/api/pb"
	"notionboy/api/pb/model"
	"notionboy/db/ent"
)

type conversationMgr struct{}

func NewConversationService() ConversationService {
	return &conversationMgr{}
}

// ConversationService is the interface for the RPC service
type ConversationService interface {
	// CreateConversation creates a new conversation
	CreateConversation(ctx context.Context, acc *ent.Account, id, instruction, title string) (*ConversationDTO, error)
	// UpdateConversation updates a conversation
	UpdateConversation(ctx context.Context, acc *ent.Account, Id, instruction, title string) (*ConversationDTO, error)

	// GetConversation gets a conversation by Id
	GetConversation(ctx context.Context, acc *ent.Account, Id string) (*ConversationDTO, error)

	// ListConversations lists all conversations
	ListConversations(ctx context.Context, acc *ent.Account, limit, offset int) ([]*ConversationDTO, error)

	// DeleteConversation deletes a conversation
	DeleteConversation(ctx context.Context, acc *ent.Account, Id string) error

	// CreateMessage creates a new conversation
	CreateConversationMessage(ctx context.Context, acc *ent.Account, conversationId, request, model string) (*ConversationMessageDTO, error)

	// CreateMessage creates a new conversation
	CreateStreamConversationMessage(ctx context.Context, acc *ent.Account, stream pb.Service_CreateMessageServer, req *model.CreateMessageRequest) error

	// GetMessage gets a conversation by Id
	GetConversationMessage(ctx context.Context, acc *ent.Account, conversationId, messageId string) (*ConversationMessageDTO, error)

	// ListMessages lists all conversations
	ListConversationMessages(ctx context.Context, acc *ent.Account, conversationId string, limit, offset int) ([]*ConversationMessageDTO, error)

	// DeleteMessage ctx context.Context, acc *ent.Account,  a conversation
	DeleteConversationMessage(ctx context.Context, acc *ent.Account, conversationId, messageId string) error
}
