package chatgpt

import (
	"context"
	"notionboy/db/ent"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"
	"sync"
)

type Chatter interface {
	// Chat
	// Params are context, parent message id and prompt
	//
	// Returns message_id, message and error
	Chat(ctx context.Context, prompt string) (string, error)
	ChatWithHistory(ctx context.Context, acc *ent.Account, prompt string) (string, error)
	ResetHistory(acc *ent.Account)
}

var (
	defaultApiClient Chatter
	once             sync.Once
)

// New create a new chatter
func New(cfg *config.ChatGPTConfig) Chatter {
	if cfg.ApiKey != "" {
		return newApiClient(cfg.ApiKey)
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
