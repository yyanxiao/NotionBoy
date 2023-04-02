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

	// set account info
	content.Account = accountInfo

	// reset chat history for ChatGPT
	if strings.ToUpper(msgText) == config.CMD_CHAT_RESET {
		ex.chatter.ResetHistory(accountInfo)
		chatPageCache.Delete(accountInfo)
		mr <- &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(config.MSG_RESET_CHATGPT_HISTORY)}
		return
	}

	if msg.MsgType != message.MsgTypeText {
		mr <- &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(config.MSG_CHAT_UNSUPPOERT)}
		return
	}

	chatter := ex.chatter
	if accountInfo.IsOpenaiAPIUser {
		chatter = chatgpt.DefaultApiClient()
	}
	resp, err := chatter.ChatWithHistory(ctx, accountInfo, strings.TrimSpace(content.Text[5:]), "")
	var chatResp string
	if err != nil {
		chatResp = err.Error()
	} else {
		chatResp = fmt.Sprintf("ChatGPT: \n%s", resp)
	}
	mr <- &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(chatResp)}
}
