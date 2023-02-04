package wxgzh

import (
	"context"
	"fmt"
	"notionboy/db/ent"
	"notionboy/internal/chatgpt"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/pkg/notion"
	"notionboy/internal/pkg/utils/cache"
	"strings"
	"time"

	"github.com/silenceper/wechat/v2/officialaccount/message"
)

var chatPageCache = &chatPageCacheClient{
	client:            cache.DefaultClient(),
	DefaultExpiration: 10 * time.Minute,
}

type chatPageCacheClient struct {
	client            *cache.Cache
	DefaultExpiration time.Duration
}

func (c *chatPageCacheClient) buildCacheKey(acc *ent.Account) string {
	return fmt.Sprintf("wxgzh:chatgpt:%s:%s", acc.UserType, acc.UserID)
}

func (c *chatPageCacheClient) Set(ctx context.Context, acc *ent.Account, notionPageID string) {
	key := c.buildCacheKey(acc)
	c.client.Set(key, notionPageID, c.DefaultExpiration)
}

func (c *chatPageCacheClient) Get(acc *ent.Account) string {
	key := c.buildCacheKey(acc)
	if v, ok := c.client.Get(key); ok {
		return v.(string)
	}
	return ""
}

func (c *chatPageCacheClient) Delete(acc *ent.Account) {
	key := c.buildCacheKey(acc)
	c.client.Delete(key)
}

func (ex *OfficialAccount) processChat(ctx context.Context, msg *message.MixMessage, content *notion.Content, mr chan *message.Reply) {
	msgText := strings.TrimSpace(msg.Content[5:])
	if len(msgText) == 0 {
		mr <- &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(config.MSG_EMPTY_MESSAGE)}
		return
	}

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

	// reset chat history for ChatGPT
	if strings.ToUpper(msgText) == config.CMD_CHAT_RESET {
		ex.chatter.ResetHistory(accountInfo)
		chatPageCache.Delete(accountInfo)
		mr <- &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(config.MSG_RESET_CHATGPT_HISTORY)}
		return
	}

	n := &notion.Notion{BearerToken: accountInfo.AccessToken, DatabaseID: accountInfo.DatabaseID}

	if msg.MsgType != message.MsgTypeText {
		mr <- &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(config.MSG_CHAT_UNSUPPOERT)}
		return
	}
	var res string

	notionPageID := chatPageCache.Get(accountInfo)
	if notionPageID == "" {
		// 创建初始 Record
		var chatPageId string
		res, chatPageId, err = n.CreateRecord(ctx, &notion.Content{
			Text: "ChatGPT 专属页面",
			Tags: []string{"chat", "ChatGPT"},
		})
		if err == nil {
			n.PageID = chatPageId
		}
		chatPageCache.Set(ctx, accountInfo, chatPageId)
	} else {
		// use previous chat page
		n.PageID = notionPageID
		res = "更新聊天记录成功，如需编辑更多，请前往 https://www.notion.so/" + notionPageID
	}
	content.ChatContent.Question = msgText
	content.IsChatContent = true

	go ex.updateChatContent(ctx, n, accountInfo, content)
	mr <- &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(res)}
}

func (ex *OfficialAccount) updateChatContent(ctx context.Context, n *notion.Notion, accountInfo *ent.Account, content *notion.Content) {
	updateLatestSchema(ctx, accountInfo, n)
	chatter := ex.chatter
	if accountInfo.IsOpenaiAPIUser {
		chatter = chatgpt.DefaultApiClient()
	}
	msg, err := chatter.ChatWithHistory(ctx, accountInfo, strings.TrimSpace(content.Text[5:]))
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
