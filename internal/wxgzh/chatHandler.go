package wxgzh

import (
	"context"
	"notionboy/db/ent"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/pkg/notion"
	"sync"

	"github.com/silenceper/wechat/v2/officialaccount/message"
)

// chatParentIdMap use to store chat parent, for context
var chatParentIdMap sync.Map

func (ex *OfficialAccount) processChat(ctx context.Context, msg *message.MixMessage, content *notion.Content, mr chan *message.Reply) {
	accountInfo, err := dao.QueryAccountByWxUser(ctx, msg.GetOpenID())
	if err != nil {
		logger.SugaredLogger.Errorf("Query Account Error: %v", err)
		mr <- &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(config.MSG_ERROR_ACCOUNT_NOT_FOUND)}
		return
	}
	if accountInfo.ID == 0 {
		mr <- bindNotion(ctx, msg)
		return
	}
	n := &notion.Notion{BearerToken: accountInfo.AccessToken, DatabaseID: accountInfo.DatabaseID}

	if msg.MsgType != message.MsgTypeText {
		mr <- &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(config.MSG_CHAT_UNSUPPOERT)}
		return
	}
	var res string

	// 创建初始 Record
	var chatPageId string
	res, chatPageId, err = n.CreateRecord(ctx, &notion.Content{
		Text: "ChatGPT 专属页面",
		Tags: []string{"chat", "ChatGPT"},
	})
	if err == nil {
		n.PageID = chatPageId
	}
	content.ChatContent.Question = content.Text[5:]
	content.IsChatContent = true
	chatParentId, ok := chatParentIdMap.Load(accountInfo.DatabaseID)
	if !ok {
		chatParentId = ""
	}
	go ex.updateChatContent(ctx, n, accountInfo, content, chatParentId.(string))
	mr <- &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(res)}
}

func (ex *OfficialAccount) updateChatContent(ctx context.Context, n *notion.Notion, accountInfo *ent.Account, content *notion.Content, chatParentId string) {
	updateLatestSchema(ctx, accountInfo, n)
	parentMessageId, msg, err := ex.chatter.Chat(chatParentId, content.Text[5:])
	chatParentIdMap.Store(accountInfo.DatabaseID, parentMessageId)
	var chatResp string
	if err != nil {
		chatResp = err.Error()
	} else {
		chatResp = msg
	}

	content.ChatContent.Answer = chatResp
	content.ChatContent.UserID = accountInfo.NotionUserID
	logger.SugaredLogger.Debugw("Content", "content", content.Text)
	n.UpdateRecord(ctx, content)
}
