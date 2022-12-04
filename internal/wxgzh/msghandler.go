package wxgzh

import (
	"notionboy/db/ent"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/pkg/utils"
	"time"

	notion "notionboy/internal/pkg/notion"

	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

var supportMsgTypeMap map[message.MsgType]bool

func init() {
	supportMsgTypeMap = map[message.MsgType]bool{
		message.MsgTypeText:  true,
		message.MsgTypeImage: true,
		message.MsgTypeVideo: true,
		message.MsgTypeVoice: true,
	}
}

func (ex *OfficialAccount) messageHandler(c *gin.Context, msg *message.MixMessage) *message.Reply {
	if msg.Event == message.EventSubscribe {
		return bindNotion(c, msg)
	}

	if msg.Event == message.EventUnsubscribe {
		return unBindingNotion(c, msg)
	}

	userID := msg.GetOpenID()
	content := transformToNotionContent(msg)
	memCache := utils.GetCache()
	userCache := memCache.Get(userID)
	logger.SugaredLogger.Infof("UserID: %s, content: %v, msgType: %s, userCache: %s", userID, content, msg.MsgType, userCache)

	switch msg.Content {
	case config.CMD_BIND:
		return bindNotion(c, msg)
	case config.CMD_UNBIND:
		return unBindingNotion(c, msg)
	case config.CMD_HELP:
		return helpInfo(c, msg)
	case config.CMD_HELP_ZH:
		return helpInfo(c, msg)
	case config.CMD_SOS:
		return sosInfo(c, msg)
	}

	mr := make(chan *message.Reply)
	go ex.processContent(c, msg, content, mr)

	select {
	case r := <-mr:
		return r
	case <-time.After(3 * time.Second):
		logger.SugaredLogger.Warnf("Save record to Notion timeout")
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(config.MSG_PROCESSING)}
	}
}

func (ex *OfficialAccount) processContent(c *gin.Context, msg *message.MixMessage, content *notion.Content, mr chan *message.Reply) {
	accountInfo, err := dao.QueryAccountByWxUser(c, msg.GetOpenID())
	if err != nil {
		logger.SugaredLogger.Errorf("Query Account Error: %v", err)
		mr <- &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(config.MSG_ERROR_ACCOUNT_NOT_FOUND)}
		return
	}
	if accountInfo.ID == 0 {
		mr <- bindNotion(c, msg)
		return
	}
	n := &notion.Notion{BearerToken: accountInfo.AccessToken, DatabaseID: accountInfo.DatabaseID}

	// 如果是不支持的类型，直接返回不支持的错误
	if _, ok := supportMsgTypeMap[msg.MsgType]; !ok {
		mr <- &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(config.MSG_UNSUPPOERT)}
		return
	}

	// 创建初始 Record
	res, pageID, err := n.CreateRecord(c, &notion.Content{
		Text: "内容正在更新，请稍等",
	})
	if err == nil {
		n.PageID = pageID
		go ex.updateNotionContent(c, msg, n, accountInfo, content)
	}
	mr <- &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(res)}
}

func updateLatestSchema(ctx *gin.Context, accountInfo *ent.Account, notionConfig *notion.Notion) {
	// 如果不是最新的 Scheam，更新 Schema
	if !accountInfo.IsLatestSchema {
		if _, err := notion.UpdateDatabaseProperties(ctx, notionConfig); err != nil {
			logger.SugaredLogger.Errorf("UpdateDatabaseProperties error: %s", err.Error())
		}
		_ = dao.UpdateIsLatestSchema(ctx, accountInfo.DatabaseID, true)
	}
}

func (ex *OfficialAccount) updateNotionContent(ctx *gin.Context, msg *message.MixMessage, n *notion.Notion, accountInfo *ent.Account, content *notion.Content) {
	updateLatestSchema(ctx, accountInfo, n)
	content.Process(ctx)
	switch msg.MsgType {
	case message.MsgTypeText:
		// 保存文本信息到 Notion
		n.UpdateRecord(ctx, content)
	case message.MsgTypeImage, message.MsgTypeVideo, message.MsgTypeVoice:
		// 保存媒体信息到 Notion
		media := NewMedia(ex.officialAccount.GetContext())
		getMediaResp, err := media.getMedia(ctx, msg.MediaID, accountInfo.DatabaseID)
		if err != nil {
			// todo
			return
		}
		content.Media = notion.MediaContent{
			URL:  getMediaResp.R2URL,
			Type: getMediaResp.ContentType,
		}
		content.IsMedia = true
		n.UpdateRecord(ctx, content)
	}
}
