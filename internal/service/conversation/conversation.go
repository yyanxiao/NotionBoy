package conversation

import (
	"context"
	"fmt"
	"notionboy/db/ent"
	"notionboy/db/ent/quota"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/pkg/utils/cache"
	"time"

	"github.com/google/uuid"
	gogpt "github.com/sashabaranov/go-gpt3"
)

const DEFAULT_INSTRUCTION = `
You are ChatGPT, a large language model trained by OpenAI. You Respond as concisely as possible for each response. It's essential to respond concisely and use the same language as the user, Please keep this in mind.
`
const DEFAULT_TITLE = "ChatGPT"
const DEFAULT_MESSAGE_LIMIT = 10

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
	}

	// incr quota
	h.Quota.DailyUsed += 1
	h.Quota.MonthlyUsed += 1

	h.append(msg)
	h.saveToCache()
	return h.saveMessageToDB(msg, resp.Usage.TotalTokens)
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
		})
	}
	h.saveToCache()
	return nil
}

func (h *History) getQuotaFromDB() error {
	qt, err := loadQuota(h.Ctx, h.Account)
	if err != nil {
		logger.SugaredLogger.Errorw("loadQuota", "err", err)
		return err
	}
	h.Quota = qt

	h.saveToCache()
	return nil
}

func (h *History) saveMessageToDB(message *Message, tokenUsage int) (*ent.ConversationMessage, error) {
	conversationMessage := &ent.ConversationMessage{
		UUID:           uuid.New(),
		ConversationID: h.ConversationId,
		UserID:         h.Account.UUID,
		Request:        message.Request,
		Response:       message.Response,
		TokenUsage:     tokenUsage,
	}
	// save quota
	err := dao.IncrDailyQuota(h.Ctx, h.Account.ID, quota.CategoryChatgpt)
	if err != nil {
		logger.SugaredLogger.Errorw("incrDailyQuota", "err", err)
		return nil, err
	}
	return dao.SaveConversationMessage(h.Ctx, conversationMessage)
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
