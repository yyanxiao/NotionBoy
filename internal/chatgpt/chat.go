package chatgpt

import (
	"context"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"
	"sync"
)

type Chatter interface {
	// Chat
	// Params are context, parent message id and prompt
	//
	// Returns message_id, message and error
	Chat(ctx context.Context, parentMessageId, prompt string) (string, string, error)
}

var (
	defaultApiClient     Chatter
	defaultReverseClient Chatter
	once                 sync.Once
)

// New create a new chatter
func New(cfg *config.ChatGPTConfig) Chatter {
	if cfg.ApiKey != "" {
		return newApiClient(cfg.ApiKey)
	}
	if cfg.SessionToken != "" || (cfg.User != "" && cfg.Pass != "") {
		return newReverseClient(cfg.SessionToken, cfg.User, cfg.Pass)
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

// DefaultReverseClient creates a new client from config.ChatGPTConfig
func DefaultReverseClient() Chatter {
	once.Do(func() {
		defaultReverseClient = New(&config.ChatGPTConfig{
			SessionToken: config.GetConfig().ChatGPT.SessionToken,
			User:         config.GetConfig().ChatGPT.User,
			Pass:         config.GetConfig().ChatGPT.Pass,
		})
	})
	return defaultReverseClient
}
