package chatgpt

import (
	"context"
	"fmt"
	"notionboy/db/ent"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/pkg/utils/cache"
	"notionboy/internal/service/conversation"
	"sync"
	"time"
)

type Chatter interface {
	ChatWithHistory(ctx context.Context, acc *ent.Account, prompt string) (string, error)
	ResetHistory(acc *ent.Account)
}

var (
	defaultApiClient Chatter
	once             sync.Once
	cacheClient      = cache.DefaultClient()
	svc              = conversation.NewConversationService()
)

type History struct {
	Account        *ent.Account
	ConversationID string
}

func (h *History) buildCahceKey() string {
	if h.Account == nil {
		return ""
	}
	return fmt.Sprintf("chatgpt:%s:%s", h.Account.UserType, h.Account.UUID.String())
}

func (h *History) Load() {
	key := h.buildCahceKey()
	if key == "" {
		return
	}

	his, ok := cacheClient.Get(h.buildCahceKey())
	if !ok {
		logger.SugaredLogger.Debugw("No history found", "key", key)
		return
	}
	h.ConversationID = his.(*History).ConversationID
}

func (h *History) Save() {
	cacheClient.Set(h.buildCahceKey(), h, 24*time.Hour)
}

// New create a new chatter
func New(cfg *config.ChatGPTConfig) Chatter {
	if cfg.ApiKey != "" {
		return &chatMgr{
			Client: conversation.NewApiClient(cfg.ApiKey),
		}
	}
	logger.SugaredLogger.Error("Invalid configuration")
	return nil
}

// DefaultApiClient from Config
func DefaultApiClient() Chatter {
	once.Do(func() {
		defaultApiClient = New(&config.ChatGPTConfig{
			ApiKey: config.GetConfig().ChatGPT.ApiKey,
		})
	})
	return defaultApiClient
}

type chatMgr struct {
	Client *conversation.ConversationClient
}

func (m *chatMgr) ChatWithHistory(ctx context.Context, acc *ent.Account, prompt string) (string, error) {
	history := &History{
		Account: acc,
	}
	history.Load()

	if history.ConversationID == "" {
		cvs, err := svc.CreateConversation(ctx, acc, "", "")
		if err != nil {
			return "", err
		}
		history.ConversationID = cvs.ID
	}

	logger.SugaredLogger.Debugw("Get prompt message for api client", "prompt", prompt, "conversationId", history.ConversationID)

	msg, err := m.Client.ChatWithHistory(ctx, acc, "", history.ConversationID, prompt)
	if err != nil {
		return "", err
	}

	history.Save()
	return msg.Response, nil
}

func (m *chatMgr) ResetHistory(acc *ent.Account) {
	cacheClient.Delete((&History{Account: acc}).buildCahceKey())
}
