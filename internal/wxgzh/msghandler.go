package wxgzh

import (
	"time"

	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/db"
	notion "notionboy/internal/pkg/notion"
	"notionboy/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	log "github.com/sirupsen/logrus"
)

func (ex *OfficialAccount) messageHandler(c *gin.Context, msg *message.MixMessage) *message.Reply {
	if msg.MsgType == message.MsgType(message.EventSubscribe) {
		return bindNotion(c, msg)
	}

	if msg.MsgType == message.MsgType(message.EventUnsubscribe) {
		return unBindingNotion(c, msg)
	}

	userID := msg.GetOpenID()
	content := transformToNotionContent(msg)
	memCache := utils.GetCache()
	userCache := memCache.Get(userID)
	log.Infof("UserID: %s, content: %v, msgType: %s, userCache: %s", userID, content, msg.MsgType, userCache)

	switch msg.Content {
	case config.CMD_BIND:
		return bindNotion(c, msg)
	case config.CMD_UNBIND:
		return unBindingNotion(c, msg)
	case config.CMD_HELP:
		return helpInfo(c, msg)
	case config.CMD_HELP_ZH:
		return helpInfo(c, msg)
	}

	// 获取用户信息
	accountInfo := db.QueryAccountByWxUser(msg.GetOpenID())
	if accountInfo.ID == 0 {
		return bindNotion(c, msg)
	}

	// 保存内容到 Notion
	var ch chan string
	go func(ch chan string) {
		notionConfig := &notion.NotionConfig{BearerToken: accountInfo.AccessToken, DatabaseID: accountInfo.DatabaseID}
		// 如果不是最新的 Scheam，更新 Schema
		if !accountInfo.IsLatestSchema {
			if _, err := notion.UpdateDatabaseProperties(c, notionConfig); err != nil {
				log.Errorf("UpdateDatabaseProperties error: %s", err.Error())
			}
			db.UpdateIsLatestSchema(accountInfo.DatabaseID, true)
		}

		switch msg.MsgType {
		case message.MsgTypeText:
			// 保存文本信息到 Notion
			res, _ := notion.CreateNewRecord(c, notionConfig, content)
			ch <- res
		case message.MsgTypeImage, message.MsgTypeVideo, message.MsgTypeVoice:
			// 保存媒体信息到 Notion
			media := NewMedia(ex.officialAccount.GetContext())
			getMediaResp, err := media.getMedia(c, msg.MediaID, accountInfo.DatabaseID)
			if err != nil {
				ch <- err.Error()
			}
			res, _ := notion.CreateNewMediaRecord(c, notionConfig, getMediaResp.R2URL, getMediaResp.ContentType)
			ch <- res
		default:
			ch <- config.MSG_UNSUPPOERT
		}
	}(ch)

	// 设置超时时间
	select {
	case s := <-ch:
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(s)}
	case <-time.After(time.Second * 3):
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(config.MSG_PROCESSING)}
	}
}
