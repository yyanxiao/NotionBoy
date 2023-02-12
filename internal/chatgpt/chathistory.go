package chatgpt

import (
	"context"
	"fmt"
	"notionboy/db/ent"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/pkg/utils/cache"
	"strings"
	"time"

	"github.com/google/uuid"
)

const DEFAULT_PROMPT = `
You are ChatGPT, a large language model trained by OpenAI. You Respond as concisely as possible for each response.
It's essential to respond concisely and use the same language as the user, Please keep this in mind.

User: Hello
ChatGPT: Hi, how are you? <|im_end|>

User: 您好
ChatGPT: 您好, 有什么可以帮您? <|im_end|>

`

type ConversationWithHistory struct {
	Ctx             context.Context `json:"ctx"`
	Acc             *ent.Account    `json:"acc"`
	ConversationID  string          `json:"conversation_id"`
	ConversationIdx int             `json:"conversation_idx"`
	MessageID       string          `json:"message_id"`
	MessageIdx      int             `json:"message_idx"`
	Prompt          string          `json:"prompt"`
	Response        string          `json:"response"`
	History         string          `json:"history"`
}

var chcheClient = cache.DefaultClient()

func buildChatHitoryKey(acc *ent.Account) string {
	return fmt.Sprintf("chatgpt:%s:%s", acc.UserType, acc.UserID)
}

func NewConversationWithHistory(ctx context.Context, acc *ent.Account) *ConversationWithHistory {
	maxConversationIdx, err := dao.QueryMaxConversationIdx(ctx, acc.ID)
	if err != nil {
		maxConversationIdx = 0
	}
	return &ConversationWithHistory{
		Ctx:             ctx,
		Acc:             acc,
		ConversationID:  uuid.New().String(),
		ConversationIdx: maxConversationIdx + 1,
		MessageIdx:      1,
	}
}

func (c *ConversationWithHistory) BuildPromptWithHistory() string {
	sb := strings.Builder{}
	sb.Write([]byte(DEFAULT_PROMPT))
	sb.Write([]byte(fmt.Sprintf("Current date: %s\n\n", time.Now().Format(time.RFC3339))))
	sb.Write([]byte(fmt.Sprintf("Chat history: %s\n\n", c.History)))
	sb.Write([]byte(fmt.Sprintf("User: %s\n", c.Prompt)))
	return sb.String()
}

func (c *ConversationWithHistory) BuildNewPrompt() string {
	newPromot := c.BuildPromptWithHistory()
	if len(newPromot) > 2048 {
		c.SummaryHistory()
		newPromot = c.BuildPromptWithHistory()
	}
	return newPromot
}

func (c *ConversationWithHistory) SummaryHistory() {
	sb := strings.Builder{}
	sb.Write([]byte("Please summary the chat history, make it short and concise, and don't be verbose, it must be less than 1024 characters."))
	sb.Write([]byte(fmt.Sprintf("-------\nChat History: %s\n", c.History)))
	prompt := sb.String()
	resp, err := defaultApiClient.Chat(c.Ctx, prompt)
	if err != nil {
		logger.SugaredLogger.Errorw("Summary history error", "error", err)
		return
	}
	c.History = ProcessResponse(resp)
}

func (c *ConversationWithHistory) CacheHistory() {
	if c.Acc == nil {
		return
	}
	key := buildChatHitoryKey(c.Acc)
	newHistory := fmt.Sprintf("%s\nUser: %s\nChatGPT: %s", c.History, c.Prompt, strings.Replace(c.Response, "ChatGPT: ", "", 1))
	c.History = newHistory
	chcheClient.Set(key, c, 10*time.Minute)
}

func (c *ConversationWithHistory) SaveHistory() error {
	if c.Acc == nil {
		return nil
	}
	err := dao.SaveChatHistory(c.Ctx, &ent.ChatHistory{
		UserID:          c.Acc.ID,
		ConversationID:  uuid.MustParse(c.ConversationID),
		ConversationIdx: c.ConversationIdx,
		MessageID:       c.MessageID,
		MessageIdx:      c.MessageIdx,
		Request:         c.Prompt,
		Response:        c.Response,
	})
	if err != nil {
		logger.SugaredLogger.Errorw("Save chat history error", "error", err)
		return err
	}
	return nil
}

// save chat history to cache, expire in 10 minutes
func setChatHistory(resp *ConversationWithHistory) {
	resp.CacheHistory()
}

func resetChatHistory(acc *ent.Account) {
	key := buildChatHitoryKey(acc)
	chcheClient.Delete(key)
	logger.SugaredLogger.Debugw("Reset chat history", "key", key)
}

// get chat history from cache
func getChatHistory(ctx context.Context, acc *ent.Account) *ConversationWithHistory {
	if acc == nil {
		logger.SugaredLogger.Info("Account is nil, no chat history")
		return NewConversationWithHistory(ctx, acc)
	}
	key := buildChatHitoryKey(acc)
	if v, ok := chcheClient.Get(key); ok {
		return v.(*ConversationWithHistory)
	}
	logger.SugaredLogger.Debugw("No chat history", "key", key)
	return NewConversationWithHistory(ctx, acc)
}
