package handler

import (
	"context"
	"strings"
	"time"

	"notionboy/internal/chatgpt"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"

	tele "gopkg.in/telebot.v3"
)

func OnChatGPT(c tele.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*60*time.Second)
	defer cancel()

	acc, err := queryUserAccount(ctx, c)
	if err != nil {
		logger.SugaredLogger.Errorf("Query Account Error: %v", err)
	}

	prompt := c.Message().Payload

	if strings.ToUpper(prompt) == config.CMD_CHAT_RESET {
		chatgpt.DefaultApiClient().ResetHistory(acc)
		return c.Reply(config.MSG_RESET_CHATGPT_HISTORY)
	}

	msg, err := chatgpt.DefaultApiClient().ChatWithHistory(ctx, acc, prompt, "")
	if err != nil {
		return c.Reply(err.Error())
	}
	return c.Reply(msg)
}
