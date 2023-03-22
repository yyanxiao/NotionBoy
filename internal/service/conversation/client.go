package conversation

import (
	"context"
	"errors"
	"io"
	"strings"

	"notionboy/api/pb"
	"notionboy/db/ent"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"

	gogpt "github.com/sashabaranov/go-openai"
)

const DEFAULT_MODEL = gogpt.GPT3Dot5Turbo

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
	if err != nil {
		return nil, err
	}
	if history.isRateLimit {
		return nil, errors.New(config.MSG_ERROR_QUOTA_LIMIT)
	}

	reqMsg := history.buildRequestMessages(prompt)

	req := gogpt.ChatCompletionRequest{
		Model:    DEFAULT_MODEL,
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

func (cli *ConversationClient) StreamChatWithHistory(ctx context.Context, acc *ent.Account, instruction, conversationId, prompt string, stream pb.Service_CreateMessageServer) (*ent.ConversationMessage, error) {
	logger.SugaredLogger.Debugw("Get prompt message for api client", "prompt", prompt, "conversationId", conversationId, "instruction", instruction)
	h := NewHistory(ctx, acc, conversationId, instruction)
	err := h.Load()
	if err != nil {
		return nil, err
	}
	if h.isRateLimit {
		return nil, errors.New(config.MSG_ERROR_QUOTA_LIMIT)
	}

	reqMsg := h.buildRequestMessages(prompt)

	req := gogpt.ChatCompletionRequest{
		Model:    DEFAULT_MODEL,
		Messages: reqMsg,
		Stream:   true,
	}

	conversationMessage := &ent.ConversationMessage{
		ConversationID: h.ConversationId,
		UserID:         acc.UUID,
		Request:        prompt,
	}

	streamResp, err := cli.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return conversationMessage, err
	}
	sb := strings.Builder{}
	for {
		response, err := streamResp.Recv()
		if errors.Is(err, io.EOF) {
			// todo save conversation message
			logger.SugaredLogger.Debugw("Stream chat EOF")
			conversationMessage.Response = sb.String()

			msg := &Message{
				Request:  prompt,
				Response: sb.String(),
			}

			h.Quota.DailyUsed += 1
			h.Quota.MonthlyUsed += 1

			h.append(msg)
			h.saveToCache()
			return h.saveMessageToDB(msg, 0)
		}

		if err != nil {
			logger.SugaredLogger.Errorw("Stream chat error", "error", err)
			return nil, err
		}

		msg := response.Choices[0].Delta.Content
		// logger.SugaredLogger.Debugw("Stream chat response", "response", msg)
		sb.WriteString(msg)
		conversationMessage.Response = msg
		dto := ConversationMessageDTOFromDB(conversationMessage)
		if err = stream.Send(dto.ToPB()); err != nil {
			logger.SugaredLogger.Errorw("Stream chat send error", "error", err)
			return conversationMessage, err
		}
	}
}
