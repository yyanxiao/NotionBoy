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

	"github.com/google/uuid"
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

func (cli *ConversationClient) StreamChatWithHistoryUpdate(ctx context.Context, acc *ent.Account, instruction string, req *model.UpdateMessageRequest, stream pb.Service_UpdateMessageServer) (*ent.ConversationMessage, error) {
	return streamChatWithHistory(ctx, cli, acc, instruction, req, stream)
}

func (cli *ConversationClient) StreamChatWithHistory(ctx context.Context, acc *ent.Account, instruction string, req *model.CreateMessageRequest, stream pb.Service_CreateMessageServer) (*ent.ConversationMessage, error) {
	return streamChatWithHistory(ctx, cli, acc, instruction, req, stream)
}

func streamChatWithHistory(ctx context.Context, cli *ConversationClient, acc *ent.Account, instruction string, req interface{}, s interface{}) (*ent.ConversationMessage, error) {
	// pre process
	var conversationId string
	var selectedModel string
	var temperature float32
	var maxTokens int32
	var prompt string
	var messageId uuid.UUID
	var err error

	switch r := req.(type) {
	case *model.CreateMessageRequest:
		conversationId = r.ConversationId
		selectedModel = r.Model
		temperature = r.Temperature
		maxTokens = r.MaxTokens
		prompt = r.Request
	case *model.UpdateMessageRequest:
		conversationId = r.ConversationId
		selectedModel = r.Model
		temperature = r.Temperature
		maxTokens = r.MaxTokens
		prompt = r.Request
		messageId, err = uuid.Parse(r.Id)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid request")
	}

	h := NewHistory(ctx, acc, conversationId, instruction)
	err = h.Load()
	if err != nil {
		return nil, err
	}
	if h.isRateLimit {
		return nil, errors.New(config.MSG_ERROR_QUOTA_LIMIT)
	}

	if selectedModel == "" {
		selectedModel = DEFAULT_MODEL
	}
	if temperature == 0 {
		temperature = 1
	}
	if maxTokens == 0 {
		maxTokens = 2000
	}

	// build request
	if err := h.summaryMessages(selectedModel, prompt); err != nil {
		return nil, err
	}
	reqMsg := h.buildRequestMessages(prompt)
	chatReq := openai.ChatCompletionRequest{
		Model:       selectedModel,
		Messages:    reqMsg,
		Stream:      true,
		MaxTokens:   int(maxTokens),
		Temperature: float32(temperature),
	}
	conversationMessage := &ent.ConversationMessage{
		UUID:           messageId,
		ConversationID: h.ConversationId,
		UserID:         acc.UUID,
		Request:        prompt,
		Model:          selectedModel,
	}

	// stream result
	streamResp, err := cli.CreateChatCompletionStream(ctx, chatReq)
	if err != nil {
		return conversationMessage, err
	}
	sb := strings.Builder{}
	for {
		response, err := streamResp.Recv()
		if errors.Is(err, io.EOF) {
			conversationMessage.Response = sb.String()
			msg := &Message{
				Id:       messageId,
				Request:  prompt,
				Response: sb.String(),
				Model:    selectedModel,
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
		// nolint:staticcheck
		switch stream := s.(type) {
		case pb.Service_CreateMessageServer:
			if err = stream.Send(dto.ToPB()); err != nil {
				logger.SugaredLogger.Errorw("Stream chat send error", "error", err)
				return conversationMessage, err
			}
		case pb.Service_UpdateMessageServer:
			if err = stream.Send(dto.ToPB()); err != nil {
				logger.SugaredLogger.Errorw("Stream chat send error", "error", err)
				return conversationMessage, err
			}
		default:
			return conversationMessage, errors.New("invalid stream")
		}
	}
}
