package handler

import (
	"context"
	"notionboy/internal/chatgpt"
	"time"

	tele "gopkg.in/telebot.v3"
)

func OnChatGPT(c tele.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	prompt := c.Message().Payload

	_, msg, err := chatgpt.DefaultApiClient().Chat(ctx, "", prompt)
	if err != nil {
		return c.Reply(err.Error())
	}
	return c.Reply(msg)
}
