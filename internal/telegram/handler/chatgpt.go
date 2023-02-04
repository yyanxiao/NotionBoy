package handler

import (
	"context"
	"notionboy/internal/chatgpt"
	"notionboy/internal/pkg/logger"
	"time"

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

	msg, err := chatgpt.DefaultApiClient().ChatWithHistory(ctx, acc, prompt)
	if err != nil {
		return c.Reply(err.Error())
	}
	return c.Reply(msg)
}
