package conversation

import (
	"context"
	"fmt"
	"sync"
	"time"

	"notionboy/db/ent"
	"notionboy/internal/pkg/db"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/pkg/utils/cache"

	"github.com/google/uuid"
	"github.com/pkoukk/tiktoken-go"
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

var (
	tk   *tiktoken.Tiktoken
	once sync.Once
)

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

func getTiktoken() *tiktoken.Tiktoken {
	if tk == nil {
		var err error
		once.Do(func() {
			tk, err = tiktoken.GetEncoding("cl100k_base")
			if err != nil {
				logger.SugaredLogger.Errorw("tiktoken.EncodingForModel", "error", err)
			}
		})
	}
	return tk
}

func calculateTokens(msg string) int {
	tk = getTiktoken()
	return len(tk.Encode(msg, nil, nil))
}

func (h *History) calculateTokens(m *Message) int64 {
	usage := m.Usage
	if usage == nil {
		promptTokens := 0
		completionTokens := 0
		tk = getTiktoken()
		promptTokens += calculateTokens(h.Instruction)
		for _, tm := range h.Messages[:len(h.Messages)-1] {
			promptTokens += calculateTokens(tm.Request)
			promptTokens += calculateTokens(tm.Response)
		}
		promptTokens += calculateTokens(m.Request)
		completionTokens += calculateTokens(m.Response)
		usage = &openai.Usage{
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
			TotalTokens:      promptTokens + completionTokens,
		}
	}

	totalTokens := 0

	// https://openai.com/pricing#faq-token
	switch m.Model {
	case openai.GPT3Dot5Turbo, openai.GPT3Dot5Turbo0301:
		totalTokens = usage.TotalTokens
	case openai.GPT4, openai.GPT40314:
		totalTokens = usage.PromptTokens*15 + usage.CompletionTokens*30
	case openai.GPT432K, openai.GPT432K0314:
		totalTokens = usage.PromptTokens*30 + usage.CompletionTokens*60
	}

	return int64(totalTokens)
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
	msg := &ent.ConversationMessage{
		UUID:           uuid.New(),
		ConversationID: h.ConversationId,
		UserID:         h.Account.UUID,
		Request:        message.Request,
		Response:       message.Response,
		TokenUsage:     usage,
		Model:          message.Model,
	}
	logger.SugaredLogger.Debugw("SaveConversationMessage to DB", "conversationMessage", msg)
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
		SetModel(message.Model).
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

func (h *History) buildRequestMessages(prompt string) []openai.ChatCompletionMessage {
	messages := make([]openai.ChatCompletionMessage, 0)
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    ROLE_SYSTEM,
		Content: h.Instruction,
	})
	for _, message := range h.Messages {
		messages = append(messages, message.toChatMessage()...)
	}
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    ROLE_USER,
		Content: prompt,
	})

	return messages
}

// summaryMessages summary messages to fit model limit
func (h *History) summaryMessages(model, prompt string) {
	tokens := 0
	modelLimit := 2000
	switch model {
	case openai.GPT4, openai.GPT40314:
		modelLimit = 4096
	case openai.GPT432K, openai.GPT432K0314:
		modelLimit = 8192
	}
	for _, message := range h.Messages {
		tokens += calculateTokens(message.Request)
		tokens += calculateTokens(message.Response)
	}
	tokens += calculateTokens(prompt)

	if tokens > modelLimit {
		logger.SugaredLogger.Debugw("history too long, summary message", "tokens", tokens, "modelLimit", modelLimit)
		// summary messages
		req := openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: h.buildRequestMessages("Please summary the conversation to less than 1000 tokens"),
		}

		resp, err := defaultApiClient.CreateChatCompletion(h.Ctx, req)
		if err != nil {
			logger.SugaredLogger.Errorw("summaryMessages", "err", err)
			return
		}
		logger.SugaredLogger.Infow("summaryMessages", "resp", resp)
		h.Messages = []*Message{
			{
				Request:  "The summary of our previous conversation",
				Response: getResponse(&resp),
				Model:    model,
			},
		}

		// update toke nusage
		tx, err := db.GetClient().Tx(h.Ctx)
		if err != nil {
			logger.SugaredLogger.Errorw("Init transaction for symmary message error", "err", err)
			return
		}
		err = dao.IncrUsedTokenQuota(tx.Client(), h.Ctx, h.Account.ID, int64(resp.Usage.TotalTokens))
		if err != nil {
			logger.SugaredLogger.Errorw("IncrUsedTokenQuota for symmary message error", "err", err)
			if e := tx.Rollback(); e != nil {
				return
			}
			return
		}
		if err := dao.IncrConversationUsedToken(tx.Client(), h.Ctx, h.ConversationId, int64(resp.Usage.TotalTokens)); err != nil {
			logger.SugaredLogger.Errorw("IncrConversationUsedToken for symmary message error", "err", err)
			if e := tx.Rollback(); e != nil {
				return
			}
			return
		}
		if err := tx.Commit(); err != nil {
			logger.SugaredLogger.Errorw("Commit on update token for summary message error", "err", err)
			return
		}
	}
}

func getResponse(resp *openai.ChatCompletionResponse) string {
	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content
	}
	return ""
}
