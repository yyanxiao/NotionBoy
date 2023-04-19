package conversation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"notionboy/db/ent"
	"notionboy/db/ent/conversationmessage"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/db"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/pkg/utils/cache"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
)

const DEFAULT_INSTRUCTION = `
You are ChatGPT, a large language model trained by OpenAI. Answer as concisely as possible. It's essential to use the same language as the user.
`

const (
	DEFAULT_TITLE         = "ChatGPT"
	DEFAULT_MESSAGE_LIMIT = 10
)

const (
	ROLE_USER      = "user"
	ROLE_ASSISTANT = "assistant"
	ROLE_SYSTEM    = "system"
)

var cacheClient = cache.DefaultClient()

type History struct {
	Ctx            context.Context
	Account        *ent.Account
	ConversationId uuid.UUID
	Instruction    string
	Messages       []*Message
	Quota          *ent.Quota
	isRateLimit    bool
	Summary        string // summary of the conversation for reduce the size of the prompt
}

type Message struct {
	Id       uuid.UUID
	Request  string
	Response string
	Model    string
	Usage    *openai.Usage
}

func (m *Message) toChatMessage() []openai.ChatCompletionMessage {
	return []openai.ChatCompletionMessage{
		{
			Role:    ROLE_USER,
			Content: m.Request,
		},
		{
			Role:    ROLE_ASSISTANT,
			Content: m.Response,
		},
	}
}

func (h *History) calculateTokens(m *Message) int64 {
	usage := m.Usage
	if usage != nil {
		tokens := calculateTotalTokensRequired(usage.PromptTokens, usage.CompletionTokens, m.Model)
		return int64(tokens)
	}

	// calculate tokens usage before save to db
	// need delete the last message from history messages as it is the new message
	// logger.SugaredLogger.Debugw("Calculate tokens usage", "historyMessages", h.Messages, "prompt", m.Request, "model", m.Model, "response", m.Response)
	messages := h.Messages[:]
	if len(messages) > 0 {
		messages = messages[:len(messages)-1]
	}
	tokens := calculateTotalTokensForMessages(messages, m.Model, h.Instruction, m.Request, m.Response)
	return int64(tokens)
}

func NewHistory(ctx context.Context, acc *ent.Account, conversationId, instruction string) *History {
	if instruction == "" {
		instruction = DEFAULT_INSTRUCTION
	}
	return &History{
		Ctx:            ctx,
		Account:        acc,
		ConversationId: uuid.MustParse(conversationId),
		Instruction:    instruction,
		Messages:       []*Message{},
	}
}

func (h *History) buildCacheKey() string {
	return fmt.Sprintf("chatgpt:conversation:%s", h.ConversationId.String())
}

func (h *History) saveToCache() {
	cacheClient.Set(h.buildCacheKey(), h, 24*time.Hour)
}

func (h *History) Load() error {
	var err error
	h.getFromCache()
	if len(h.Messages) == 0 {
		if err := h.getMessagesFromDB(); err != nil {
			return err
		}
	}
	if h.Quota == nil {
		if err := h.getQuotaFromDB(); err != nil {
			return err
		}
	}

	h.isRateLimit = checkRateLimit(h.Account, h.Quota)
	return err
}

func (h *History) Save(req *openai.ChatCompletionMessage, resp *openai.ChatCompletionResponse) (*ent.ConversationMessage, error) {
	msg := &Message{
		Request:  req.Content,
		Response: getResponse(resp),
		Model:    resp.Model,
		Usage:    &resp.Usage,
	}

	h.append(msg)
	h.saveToCache()
	return h.saveMessageToDB(msg)
}

func (h *History) getFromCache() {
	history, ok := cacheClient.Get(h.buildCacheKey())
	if ok {
		history := history.(*History)
		h.Account = history.Account
		h.ConversationId = history.ConversationId
		h.Instruction = history.Instruction
		h.Messages = history.Messages
		h.Quota = history.Quota
		h.isRateLimit = history.isRateLimit
	}
}

func (h *History) getMessagesFromDB() error {
	messages, err := dao.ListConversationMessages(h.Ctx, h.ConversationId, DEFAULT_MESSAGE_LIMIT, 0)
	if err != nil {
		logger.SugaredLogger.Errorw("getConversationHistoryFromDB", "err", err)
		return err
	}
	for _, message := range messages {
		h.append(&Message{
			Request:  message.Request,
			Response: message.Response,
			Model:    message.Model,
		})
	}
	h.saveToCache()
	return nil
}

func (h *History) getQuotaFromDB() error {
	qt, err := dao.QueryQuota(h.Ctx, h.Account.ID)
	if err != nil {
		logger.SugaredLogger.Errorw("QueryQuota", "err", err)
		return err
	}
	h.Quota = qt

	h.saveToCache()
	return nil
}

func (h *History) saveMessageToDB(message *Message) (*ent.ConversationMessage, error) {
	// usage should contains all tokens current message add message history
	usage := h.calculateTokens(message)
	var messageId uuid.UUID
	if message.Id != uuid.Nil {
		messageId = message.Id
	} else {
		messageId = uuid.New()
	}

	tx, err := db.GetClient().Tx(h.Ctx)
	if err != nil {
		logger.SugaredLogger.Errorw("SaveConversationMessage", "err", err)
		return nil, err
	}
	// save quota
	err = dao.IncrUsedTokenQuota(tx.Client(), h.Ctx, h.Account.ID, usage)

	if err != nil {
		logger.SugaredLogger.Errorw("IncrUsedTokenQuota", "err", err)
		if e := tx.Rollback(); e != nil {
			return nil, e
		}
		return nil, err
	}

	if err := dao.IncrConversationUsedToken(tx.Client(), h.Ctx, h.ConversationId, usage); err != nil {
		logger.SugaredLogger.Errorw("IncrConversationUsedToken", "err", err)
		if e := tx.Rollback(); e != nil {
			return nil, e
		}
		return nil, err
	}

	var msg *ent.ConversationMessage

	if message.Id != uuid.Nil {
		// update existing
		msg, err = tx.ConversationMessage.Query().Where(conversationmessage.UUIDEQ(message.Id)).Only(h.Ctx)
		if err != nil {
			logger.SugaredLogger.Errorw("Update ConversationMessage error", "err", err, "message", message)
			if e := tx.Rollback(); e != nil {
				return nil, e
			}
			return nil, err
		}

		if err := tx.ConversationMessage.Update().
			SetRequest(message.Request).
			SetResponse(message.Response).
			SetTokenUsage(usage).
			SetModel(message.Model).
			Where(conversationmessage.UUIDEQ(message.Id)).
			Exec(h.Ctx); err != nil {
			if err != nil {
				logger.SugaredLogger.Errorw("Update ConversationMessage error", "err", err)
				if e := tx.Rollback(); e != nil {
					return nil, e
				}
				return nil, err
			}
		}
		// delete all message that created after this message
		if _, err := tx.ConversationMessage.Delete().
			Where(conversationmessage.ConversationIDEQ(h.ConversationId), conversationmessage.UserIDEQ(msg.UserID), conversationmessage.CreatedAtGT(msg.CreatedAt)).
			Exec(h.Ctx); err != nil {
			if err != nil {
				logger.SugaredLogger.Errorw("Delete ConversationMessage error", "err", err)
				if e := tx.Rollback(); e != nil {
					return nil, e
				}
				return nil, err
			}
		}
	} else {
		// create new
		msg = &ent.ConversationMessage{
			UUID:           messageId,
			ConversationID: h.ConversationId,
			UserID:         h.Account.UUID,
			Request:        message.Request,
			Response:       message.Response,
			TokenUsage:     usage,
			Model:          message.Model,
		}
		if err := tx.ConversationMessage.Create().
			SetUserID(msg.UserID).
			SetConversationID(msg.ConversationID).
			SetRequest(msg.Request).
			SetResponse(msg.Response).
			SetTokenUsage(msg.TokenUsage).
			SetModel(message.Model).
			SetUUID(msg.UUID).Exec(h.Ctx); err != nil {
			if err != nil {
				logger.SugaredLogger.Errorw("SaveConversationMessage", "err", err)
				if e := tx.Rollback(); e != nil {
					return nil, e
				}
				return nil, err
			}
		}
	}
	// return the message
	msg, err = tx.ConversationMessage.Query().Where(conversationmessage.UUIDEQ(messageId)).Only(h.Ctx)
	if err != nil {
		logger.SugaredLogger.Errorw("Update ConversationMessage error", "err", err)
		if e := tx.Rollback(); e != nil {
			return nil, e
		}
		return nil, err
	}
	return msg, tx.Commit()
}

func (h *History) append(message *Message) {
	h.Messages = append(h.Messages, message)
}

func (h *History) buildRequestMessages(prompt string) ([]openai.ChatCompletionMessage, int) {
	messages := make([]openai.ChatCompletionMessage, 0)
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    ROLE_SYSTEM,
		Content: h.Instruction,
	})
	for _, message := range h.Messages {
		messages = append(messages, message.toChatMessage()...)
	}
	latestMessage := openai.ChatCompletionMessage{
		Role:    ROLE_USER,
		Content: prompt,
	}
	messages = append(messages, latestMessage)
	// messages already contain instruction and prompt, so we only need to calculate the total tokens for messages
	promptTokens := calculateTotalTokens(messages, openai.GPT3Dot5Turbo, "", "", "")
	return messages, promptTokens
}

// summaryMessages summary messages to fit model limit
func (h *History) summaryMessages(model, prompt string) error {
	logger.SugaredLogger.Info("history too long, summary message")
	// summary messages
	maxSummaryTokens := 1000
	symmaryTemperture := 0.5
	req, err := buildChatCompletionRequest(h.Messages, prompt, model, h.Instruction, maxSummaryTokens, float32(symmaryTemperture), false, h.Quota)
	if err != nil {
		logger.SugaredLogger.Errorw("summaryMessages", "err", err)
		return err
	}
	resp, err := defaultApiClient.CreateChatCompletion(h.Ctx, *req)
	if err != nil {
		logger.SugaredLogger.Errorw("summaryMessages", "err", err)
		return err
	}
	// update cache
	logger.SugaredLogger.Infow("summaryMessages", "resp", resp)
	h.Messages = []*Message{
		{
			Request:  "The summary of our previous conversation",
			Response: getResponse(&resp),
			Model:    model,
		},
	}

	// update toke nusage in db
	tx, err := db.GetClient().Tx(h.Ctx)
	if err != nil {
		logger.SugaredLogger.Errorw("Init transaction for symmary message error", "err", err)
		return err
	}
	err = dao.IncrUsedTokenQuota(tx.Client(), h.Ctx, h.Account.ID, int64(resp.Usage.TotalTokens))
	if err != nil {
		logger.SugaredLogger.Errorw("IncrUsedTokenQuota for symmary message error", "err", err)
		if e := tx.Rollback(); e != nil {
			return err
		}
		return err
	}
	if err := dao.IncrConversationUsedToken(tx.Client(), h.Ctx, h.ConversationId, int64(resp.Usage.TotalTokens)); err != nil {
		logger.SugaredLogger.Errorw("IncrConversationUsedToken for symmary message error", "err", err)
		if e := tx.Rollback(); e != nil {
			return err
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		logger.SugaredLogger.Errorw("Commit on update token for summary message error", "err", err)
		return err
	}
	return nil
}

func getResponse(resp *openai.ChatCompletionResponse) string {
	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content
	}
	return ""
}

func buildChatCompletionRequest(messages []*Message, model, instruction, prompt string, maxTokens int, temperature float32, isStream bool, quota *ent.Quota) (*openai.ChatCompletionRequest, error) {
	promptTokens := calculateTotalTokensForMessages(messages, openai.GPT3Dot5Turbo, instruction, prompt, "")
	if (promptTokens + 1000) > int(quota.Token-quota.TokenUsed) {
		logger.SugaredLogger.Debugw("history too long, summary message", "tokens", promptTokens)
		return nil, errors.New(config.MSG_ERROR_QUOTA_NOT_ENOUGH)
	}
	reqMsg := make([]openai.ChatCompletionMessage, 0)
	chatReq := &openai.ChatCompletionRequest{
		Model:       model,
		Messages:    reqMsg,
		Stream:      isStream,
		MaxTokens:   calculateMaxReturnTokens(promptTokens, maxTokens, model),
		Temperature: temperature,
	}
	return chatReq, nil
}
