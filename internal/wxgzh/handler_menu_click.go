package wxgzh

import (
	"context"
	"notionboy/internal/pkg/logger"

	"github.com/silenceper/wechat/v2/officialaccount/message"
)

const (
	helpNote = `直接在公众号回复任何内容, NotionBoy 会自动将内容保存到 Notion 中，例如: 「我是一条笔记」
1. 支持添加标签，例如: 「#NotionBoy #笔记 我是一条笔记」
2. 支持添加链接，例如: 「https://mp.weixin.qq.com/s/ib7HrRMIXwZjJyYFOwBQrw #NotionBoy #笔记 我是一条笔记」
3. 支持通过换行分割标题和内容，例如: 「我是标题
我是内容」	保存到 Notion 时 Title 为我是标题, Content 为我是内容
`
	helpFultext = `在公众号保存链接的时候，加上「#全文」这个标签, 可以剪辑文章全文到 Notion 中, 例如:

#全文 https://mp.weixin.qq.com/s/ib7HrRMIXwZjJyYFOwBQrw #NotionBoy
`
	helpZlib = `回复 「/zlib 书名或者作者」可以获取图书的下载信息，例如

/zlib 如何阅读一本书
`
)

func handleMenuClick(c context.Context, msg *message.MixMessage) *message.Reply {
	key := msg.EventKey
	logger.SugaredLogger.Debugw("handleMenuClick", "key", key)

	switch key {
	case BtnBind.String():
		return bindNotion(c, msg)
	case BtnUnbind.String():
		return unBindingNotion(c, msg)
	case BtnMagicCode.String():
		return magicCode(c, msg)
	case BtnhelpSOS.String():
		return sosInfo(c, msg)
	case BtnWhoAMI.String():
		return whoAMI(c, msg)
	case BtnHelpNote.String():
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(helpNote)}
	case BtnHelpFulltext.String():
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(helpFultext)}
	case BtnHelpZlib.String():
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(helpZlib)}
	}
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText("Unknown menu click")}
}
