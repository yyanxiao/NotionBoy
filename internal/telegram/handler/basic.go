package handler

import (
	"context"
	"fmt"
	"notionboy/db/ent/account"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/pkg/notion"
	"strconv"

	tele "gopkg.in/telebot.v3"
)

const HELP_MSG = `
这些命令和基本操作描述的是通过 NotionBoy 将内容保存到 Notion 中的功能。
- /start or /help 命令获取 NotionBoy 的基础功能介绍，可以帮助用户了解 NotionBoy 的功能
- /bind 命令可以用于绑定 Notion 账户，使 NotionBoy 能够访问 Notion 中的内容。
- /unbind 命令可以用于解绑 Notion 账户，使 NotionBoy 不再能够访问 Notion 中的内容。
- /chat 命令可以与 ChatGPT 畅聊，ChatGPT 是一种自然语言生成模型，能够通过对话方式回答用户的问题。


基本操作
- 发送任意文字、图片或者视频到 NotionBoy 时, NotionBoy 会将内容保存到 Notion 中
- 如果发送到内容中包含 # 开头的内容，会被自动识别成标签，并在 Notion 中添加这个标签
- 如果发送的内容中包含 #全文和一个 URL，则会自动保存此 URL 的网页截图
- 如果发送的内容包含 #pdf 和一个 URL，则会自动将此 URL 的网页保存为 PDF
`

func OnStart(c tele.Context) error {
	return c.Send(HELP_MSG)
}

func OnBind(c tele.Context) error {
	sender := c.Sender()
	if sender == nil {
		return c.Reply("User do not exist")
	}
	stage := fmt.Sprintf("%s:%d", account.UserTypeTelegram, sender.ID)
	oauthMgr := notion.GetOauthManager()
	url := oauthMgr.OAuthURL(stage)
	logger.SugaredLogger.Info("OAuthURL: ", url)
	text := config.MSG_BINDING + url
	return c.Reply(text)
}

func OnUnbind(c tele.Context) error {
	sender := c.Sender()
	if sender == nil {
		return c.Reply("User do not exist")
	}

	if err := dao.DeleteAccount(context.Background(), account.UserTypeTelegram, strconv.FormatInt(sender.ID, 10)); err != nil {
		return c.Reply(config.MSG_UNBIND_FAILED + err.Error())
	}
	return c.Reply(config.MSG_UNBIND_SUCCESS)
}
