package wxgzh

import (
	"notionboy/internal/pkg/db"
	notion "notionboy/internal/pkg/notion"
	"notionboy/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	log "github.com/sirupsen/logrus"
)

func messageHandler(c *gin.Context, msg *message.MixMessage) *message.Reply {

	if msg.MsgType == message.MsgType(message.EventSubscribe) {
		return bindNotion(c, msg)
	}

	userID := msg.GetOpenID()
	content := transformToNotionContent(msg)
	memCache := utils.GetCache()
	userCache := memCache.Get(userID)
	log.Infof("UserID: %s, content: %s, msgType: %s, userCache: %s", userID, content, msg.MsgType, userCache)

	if msg.Content == "绑定" {
		return bindNotion(c, msg)
	} else if msg.Content == "解绑" {
		return unBindingNotion(c, msg)
	}

	// 获取用户信息
	accountInfo := db.QueryAccountByWxUser(msg.GetOpenID())
	if accountInfo.ID == 0 {
		return bindNotion(c, msg)
	}

	res, _ := notion.CreateNewRecord(c, &notion.NotionConfig{BearerToken: accountInfo.AccessToken, DatabaseID: accountInfo.DatabaseID}, content)
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(res)}
}
