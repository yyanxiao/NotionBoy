package chatgpt

import (
	"context"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"
	"sync"
)

type Chatter interface {
	Chat(ctx context.Context, parentMessageId, prompt string) (string, string, error)
}

var (
	defaultApiClient     Chatter
	defaultReverseClient Chatter
)

// New create a new chatter
func New(cfg *config.ChatGPTConfig) Chatter {
	if cfg.ApiKey != "" {
		return newApiClient(cfg.ApiKey)
	}
	if cfg.SessionToken != "" || (cfg.User != "" && cfg.Pass != "") {
		return newReverseClient(cfg.SessionToken)
	}
	logger.SugaredLogger.Error("Invalid configuration")
	return nil
}

// DefaultApiClient from Config
func DefaultApiClient() Chatter {
	var once sync.Once
	once.Do(func() {
		defaultApiClient = New(&config.ChatGPTConfig{
			ApiKey: config.GetConfig().ChatGPT.ApiKey,
		})
	})
	return defaultApiClient
}

// DefaultReverseClient creates a new client from config.ChatGPTConfig
func DefaultReverseClient() Chatter {
	var once sync.Once
	once.Do(func() {
		defaultReverseClient = New(&config.ChatGPTConfig{
			SessionToken: config.GetConfig().ChatGPT.SessionToken,
			User:         config.GetConfig().ChatGPT.User,
			Pass:         config.GetConfig().ChatGPT.Pass,
		})
	})
	return defaultReverseClient
}
