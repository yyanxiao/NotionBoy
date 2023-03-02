package conversation

import (
	"context"
	"errors"
	"notionboy/db/ent"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"

	gogpt "github.com/sashabaranov/go-gpt3"
)

type ConversationClient struct {
	*gogpt.Client
}

func newApiClient(apiKey string) *ConversationClient {
	client := &ConversationClient{
		Client: gogpt.NewClient(apiKey),
	}
	return client
}

var defaultApiClient = newApiClient(config.GetConfig().ChatGPT.ApiKey)

func NewApiClient(apiKey string) *ConversationClient {
	if apiKey == "" {
		return defaultApiClient
	}
	return newApiClient(apiKey)
}

func (cli *ConversationClient) ChatWithHistory(ctx context.Context, acc *ent.Account, instruction, conversationId, prompt string) (*ent.ConversationMessage, error) {
	logger.SugaredLogger.Debugw("Get prompt message for api client", "prompt", prompt, "conversationId", conversationId, "instruction", instruction)
	history := NewHistory(ctx, acc, conversationId, instruction)
	err := history.Load()
	// _, isRateLimit, err := checkRateLimit(ctx, acc)
	if err != nil {
		return nil, err
	}
	if history.isRateLimit {
		return nil, errors.New(config.MSG_ERROR_QUOTA_LIMIT)
	}

	reqMsg := history.buildRequestMessages(prompt)

	req := gogpt.ChatCompletionRequest{
		Model:    gogpt.GPT3Dot5Turbo0301,
		Messages: reqMsg,
	}

	resp, err := cli.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, err
	}
	msg, err := history.Save(&reqMsg[len(reqMsg)-1], &resp)
	if err != nil {
		logger.SugaredLogger.Errorw("Save conversation message error", "error", err)
		return nil, err
	}
	return msg, nil
}
