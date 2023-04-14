package conversation

import (
	"context"
	"errors"
	"io"
	"strings"

	"notionboy/api/pb"
	"notionboy/api/pb/model"
	"notionboy/db/ent"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"

	"github.com/sashabaranov/go-openai"
)

const DEFAULT_MODEL = openai.GPT3Dot5Turbo

type ConversationClient struct {
	*openai.Client
}

func newApiClient(apiKey string) *ConversationClient {
	client := &ConversationClient{
		Client: openai.NewClient(apiKey),
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

func (cli *ConversationClient) ChatWithHistory(ctx context.Context, acc *ent.Account, instruction, conversationId, prompt, model string) (*ent.ConversationMessage, error) {
	logger.SugaredLogger.Debugw("Get prompt message for api client", "prompt", prompt, "conversationId", conversationId, "instruction", instruction)
	history := NewHistory(ctx, acc, conversationId, instruction)
	err := history.Load()
	if err != nil {
		return nil, err
	}
	if history.isRateLimit {
		return nil, errors.New(config.MSG_ERROR_QUOTA_LIMIT)
	}

	selectModel := model
	if model == "" {
		selectModel = DEFAULT_MODEL
	}
	if err := history.summaryMessages(selectModel, prompt); err != nil {
		return nil, err
	}
	reqMsg := history.buildRequestMessages(prompt)

	req := openai.ChatCompletionRequest{
		Model:     selectModel,
		Messages:  reqMsg,
		MaxTokens: 2000,
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

func (cli *ConversationClient) StreamChatWithHistory(ctx context.Context, acc *ent.Account, instruction string, req *model.CreateMessageRequest, stream pb.Service_CreateMessageServer) (*ent.ConversationMessage, error) {
	logger.SugaredLogger.Debugw("Get prompt message for api client", "req", req)
	conversationId := req.ConversationId
	h := NewHistory(ctx, acc, conversationId, instruction)
	err := h.Load()
	if err != nil {
		return nil, err
	}
	if h.isRateLimit {
		return nil, errors.New(config.MSG_ERROR_QUOTA_LIMIT)
	}

	if req.Model == "" {
		req.Model = DEFAULT_MODEL
	}
	if req.Temperature == 0 {
		req.Temperature = 1
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = 2000
	}

	if err := h.summaryMessages(req.Model, req.Request); err != nil {
		return nil, err
	}
	reqMsg := h.buildRequestMessages(req.Request)

	chatReq := openai.ChatCompletionRequest{
		Model:       req.Model,
		Messages:    reqMsg,
		Stream:      true,
		MaxTokens:   int(req.MaxTokens),
		Temperature: req.Temperature,
	}

	conversationMessage := &ent.ConversationMessage{
		ConversationID: h.ConversationId,
		UserID:         acc.UUID,
		Request:        req.Request,
		Model:          req.Model,
	}

	streamResp, err := cli.CreateChatCompletionStream(ctx, chatReq)
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
				Request:  req.Request,
				Response: sb.String(),
				Model:    req.Model,
			}
			h.append(msg)
			h.saveToCache()
			return h.saveMessageToDB(msg)
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
