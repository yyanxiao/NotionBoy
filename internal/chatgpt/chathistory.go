package chatgpt

import (
	"context"
	"fmt"
	"notionboy/db/ent"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/pkg/utils/cache"
	"strings"
	"time"
)

const DEFAULT_PROMPT = `
You are ChatGPT, a large language model trained by OpenAI. You Respond as concisely as possible for each response (e.g. don't be verbose).
It's essential to respond concisely and use the same language as the user, Please keep this in mind.

User: Hello
ChatGPT: Hi, how are you? <|im_end|>

User: 您好
ChatGPT: 您好, 有什么可以帮您? <|im_end|>


`

var chcheClient = cache.DefaultClient()

type ChatHistory struct {
	ParentMessageId string `json:"parent_message_id"`
	Prompt          string `json:"prompt"`
	Response        string `json:"response"`
	History         string `json:"history"`
}

func buildChatHitoryKey(acc *ent.Account) string {
	return fmt.Sprintf("chatgpt:%s:%s", acc.UserType, acc.UserID)
}

// save chat history to cache, expire in 10 minutes
func setChatHistory(ctx context.Context, acc *ent.Account, prompt, response string) {
	key := buildChatHitoryKey(acc)
	history := fmt.Sprintf("%s\nUser: %s\nChatGPT: %s", getChatHistory(ctx, acc), prompt, strings.Replace(response, "ChatGPT: ", "", 1))
	chcheClient.Set(key, history, 10*time.Minute)
}

func resetChatHistory(acc *ent.Account) {
	key := buildChatHitoryKey(acc)
	chcheClient.Delete(key)
	logger.SugaredLogger.Debugw("Reset chat history", "key", key)
}

// get chat history from cache
func getChatHistory(ctx context.Context, acc *ent.Account) string {
	if acc == nil {
		logger.SugaredLogger.Info("Account is nil, no chat history")
		return ""
	}
	key := buildChatHitoryKey(acc)
	if v, ok := chcheClient.Get(key); ok {
		return v.(string)
	}
	logger.SugaredLogger.Debugw("No chat history", "key", key)
	return ""
}

func builPromptdwithHistory(prompt, history string) string {
	sb := strings.Builder{}
	sb.Write([]byte(DEFAULT_PROMPT))
	sb.Write([]byte(fmt.Sprintf("Current date: %s\n\n", time.Now().Format(time.RFC3339))))
	sb.Write([]byte(fmt.Sprintf("Chat history: %s\n\n", history)))
	sb.Write([]byte(fmt.Sprintf("User: %s\n", prompt)))
	return sb.String()
}

func buildPrompt(ctx context.Context, acc *ent.Account, prompt string) string {
	history := getChatHistory(ctx, acc)
	newPromot := builPromptdwithHistory(prompt, history)
	if len(newPromot) > 2048 {
		summaryHistory(ctx, acc, history)
		history = getChatHistory(ctx, acc)
		newPromot = builPromptdwithHistory(prompt, history)
	}
	return newPromot
}

func summaryHistory(ctx context.Context, acc *ent.Account, history string) string {
	sb := strings.Builder{}
	sb.Write([]byte("Please summary the chat history, make it short and concise, and don't be verbose, it must be less than 1024 characters."))
	sb.Write([]byte(fmt.Sprintf("-------\nChat History: %s\n", history)))
	prompt := sb.String()
	resp, err := defaultApiClient.Chat(ctx, prompt)
	if err != nil {
		logger.SugaredLogger.Errorw("Summary history error", "error", err)
		return history
	}

	// save summary history
	chcheClient.Set(buildChatHitoryKey(acc), history, 10*time.Minute)

	return resp
}
