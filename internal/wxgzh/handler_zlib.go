package wxgzh

import (
	"context"
	"fmt"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/zlib"

	notion "notionboy/internal/pkg/notion"

	"github.com/silenceper/wechat/v2/officialaccount/message"
)

func (ex *OfficialAccount) processZlib(ctx context.Context, msg *message.MixMessage, content *notion.Content, mr chan *message.Reply) {
	acc, err := dao.QueryAccountByWxUser(ctx, msg.GetOpenID())
	if err != nil {
		logger.SugaredLogger.Errorf("Query Account Error: %v", err)
		mr <- &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(config.MSG_ERROR_ACCOUNT_NOT_FOUND)}
		return
	}
	if acc.ID == 0 {
		mr <- bindNotion(ctx, msg)
		return
	}
	n := &notion.Notion{BearerToken: acc.AccessToken, DatabaseID: acc.DatabaseID}

	if msg.MsgType != message.MsgTypeText {
		mr <- &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(config.MSG_ZLIB_UNSUPPOERT)}
		return
	}
	var res string
	name := content.Text[5:]
	logger.SugaredLogger.Debugw("zlib search name", "name", name)
	// 创建初始 Record
	var zlibPageId string
	res, zlibPageId, err = n.CreateRecord(ctx, &notion.Content{
		Text: "Zlib 专属页面",
	})
	if err != nil {
		mr <- &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(fmt.Sprintf("Create zlib page error: %s", err))}
		return
	}
	n.PageID = zlibPageId

	books, err := zlib.DefaultZlibClient().Search(ctx, content.Text[5:])
	if err != nil {
		mr <- &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(fmt.Sprintf("Search from zlib error: %s", err))}
		return
	}
	nContent := &notion.Content{
		Zlib:    &notion.ZlibContent{Books: books},
		Text:    name,
		Tags:    []string{"zlib", "wechat", name},
		Account: acc,
	}

	nContent.Process(ctx)
	n.UpdateRecord(ctx, nContent)

	mr <- &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(res)}
}
