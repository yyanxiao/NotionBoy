package wechat

import (
	"context"
	"notionboy/internal/chatgpt"
	"notionboy/internal/pkg/logger"

	"github.com/eatmoreapple/openwechat"
)

func chat(msg *openwechat.MessageContext) {
	if len(msg.Content) > 5 && msg.Content[:5] == "#chat" {

		logger.SugaredLogger.Infow("Receive chat message", "msg", msg.Content, "user", msg.FromUserName)
		ctx := context.Background()

		chatter := chatgpt.DefaultApiClient()
		if err := replay(msg, "收到消息，正在思考🤔"); err != nil {
			return
		}

		_, txt, err := chatter.Chat(ctx, "", msg.Content[5:])
		if err != nil {
			if err := replay(msg, err.Error()); err != nil {
				return
			}
		}
		if err := replay(msg, txt); err != nil {
			return
		}
	}
}

func replay(msg *openwechat.MessageContext, text string) error {
	if _, err := msg.ReplyText(text); err != nil {
		logger.SugaredLogger.Errorw("Send message error", "err", err, "user", msg.FromUserName)
		return err
	}
	return nil
}
