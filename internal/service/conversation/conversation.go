package conversation

import (
	"context"
	"fmt"
	"notionboy/db/ent"
	"notionboy/internal/pkg/db"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/pkg/utils/cache"
	"time"

	"github.com/google/uuid"
	gogpt "github.com/sashabaranov/go-openai"
)

const DEFAULT_INSTRUCTION = `
You are ChatGPT, a large language model trained by OpenAI. You Respond as concisely as possible for each response. It's essential to respond concisely and use the same language as the user, Please keep this in mind.
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
}

type Message struct {
	Request  string
	Response string
	Model    string
	Usage    *gogpt.Usage
}

func (m *Message) toChatMessage() []gogpt.ChatCompletionMessage {
	return []gogpt.ChatCompletionMessage{
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

func (m *Message) calculateTokens() int64 {
	// 1. get raw tokens
	// 2. calculate tokens base on model
	usage := m.Usage

	if usage == nil {
		promptTokens := len(m.Request)
		completionTokens := len(m.Response)
		usage = &gogpt.Usage{
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
			TotalTokens:      promptTokens + completionTokens,
		}
	}

	totalTokens := 0

	// https://openai.com/pricing#faq-token
	switch m.Model {
	case gogpt.GPT3Dot5Turbo, gogpt.GPT3Dot5Turbo0301:
		totalTokens = usage.TotalTokens
	case gogpt.GPT4, gogpt.GPT40314:
		totalTokens = usage.PromptTokens*15 + usage.CompletionTokens*30
	case gogpt.GPT432K, gogpt.GPT432K0314:
		totalTokens = usage.PromptTokens*30 + usage.CompletionTokens*60
	}

	return int64(totalTokens)
}

func NewHistory(ctx context.Context, acc *ent.Account, conversationId string, instruction string) *History {
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

func (h *History) Save(req *gogpt.ChatCompletionMessage, resp *gogpt.ChatCompletionResponse) (*ent.ConversationMessage, error) {
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
	usage := message.calculateTokens()

	msg := &ent.ConversationMessage{
		UUID:           uuid.New(),
		ConversationID: h.ConversationId,
		UserID:         h.Account.UUID,
		Request:        message.Request,
		Response:       message.Response,
		TokenUsage:     usage,
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

	if err := tx.ConversationMessage.Create().
		SetUserID(msg.UserID).
		SetConversationID(msg.ConversationID).
		SetRequest(msg.Request).
		SetResponse(msg.Response).
		SetTokenUsage(msg.TokenUsage).
		SetUUID(uuid.New()).Exec(h.Ctx); err != nil {
		if err != nil {
			logger.SugaredLogger.Errorw("SaveConversationMessage", "err", err)
			if e := tx.Rollback(); e != nil {
				return nil, e
			}
			return nil, err
		}
	}
	return msg, tx.Commit()
}

func (h *History) append(message *Message) {
	h.Messages = append(h.Messages, message)
}

func (h *History) buildRequestMessages(prompt string) []gogpt.ChatCompletionMessage {
	messages := make([]gogpt.ChatCompletionMessage, 0)
	messages = append(messages, gogpt.ChatCompletionMessage{
		Role:    ROLE_SYSTEM,
		Content: h.Instruction,
	})
	for _, message := range h.Messages {
		messages = append(messages, message.toChatMessage()...)
	}
	messages = append(messages, gogpt.ChatCompletionMessage{
		Role:    ROLE_USER,
		Content: prompt,
	})

	return messages
}

func getResponse(resp *gogpt.ChatCompletionResponse) string {
	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content
	}
	return ""
}
