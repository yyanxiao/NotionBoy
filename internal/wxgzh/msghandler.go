package wxgzh

import (
	"context"
	"fmt"
	"notionboy/db/ent"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/logger"
	"strconv"
	"strings"
	"time"

	notion "notionboy/internal/pkg/notion"

	"github.com/silenceper/wechat/v2/officialaccount/message"
	"golang.org/x/sync/singleflight"
)

var supportMsgTypeMap map[message.MsgType]bool

var sg singleflight.Group

func init() {
	supportMsgTypeMap = map[message.MsgType]bool{
		message.MsgTypeText:  true,
		message.MsgTypeImage: true,
		message.MsgTypeVideo: true,
		message.MsgTypeVoice: true,
	}
}

func (ex *OfficialAccount) messageHandler(ctx context.Context, msg *message.MixMessage) *message.Reply {
	switch msg.Event {
	case message.EventSubscribe:
		return helpInfo(ctx, msg)
	case message.EventUnsubscribe:
		unsubscribe(ctx, msg)
		return nil
	case message.EventScan:
		return scanQrcode(ctx, msg)
	case message.EventClick:
		return handleMenuClick(ctx, msg)
	}
	// TrimSpace will remove all space in the beginning and end of the string for matching commands
	msg.Content = strings.TrimSpace(msg.Content)
	content := transformToNotionContent(msg)

	mr := make(chan *message.Reply)
	msgID := strconv.FormatInt(msg.MsgID, 10)
	defer sg.Forget(msgID)

	isChat := func() bool {
		return strings.HasPrefix(strings.ToUpper(msg.Content), config.CMD_CHAT)
	}

	// singleflight.Group Do will process wechat retry logic
	res, _, _ := sg.Do(msgID, func() (interface{}, error) {
		cmd := strings.ToUpper(msg.Content)
		switch cmd {
		case config.CMD_BIND:
			return bindNotion(ctx, msg), nil
		case config.CMD_UNBIND:
			return unBindingNotion(ctx, msg), nil
		case config.CMD_HELP, config.CMD_HELP_ZH:
			return helpInfo(ctx, msg), nil
		case config.CMD_SOS:
			return sosInfo(ctx, msg), nil
		case config.CMD_ZLIB_NEXT:
			return searchZlibNextPage(ctx, msg), nil
		case config.CMD_ZLIB_SAVE_TO_NOTION:
			return searchZlibSaveToNotion(context.TODO(), msg), nil
		case config.CMD_UI:
			return webui(ctx, msg), nil
		case config.CMD_MAGIC_CODE:
			return magicCode(ctx, msg), nil
		case config.CMD_WHOAMI:
			return whoAMI(ctx, msg), nil
		}

		// process chatGPT
		if isChat() {
			go ex.processChat(context.TODO(), msg, content, mr)
		} else if strings.HasPrefix(strings.ToUpper(msg.Content), config.CMD_ZLIB) {
			go searchZlib(context.TODO(), msg, mr)
		} else {
			go ex.processContent(context.TODO(), msg, content, mr)
		}

		select {
		case r := <-mr:
			return r, nil
		// wechat timeout set to 13 seconds
		case <-time.After(13 * time.Second):
			if isChat() {
				text := fmt.Sprintf("ChatGPT 不能及时返回，请到 %s%s 查看结果, 网页上包含所有的聊天记录",
					config.GetConfig().Service.URL, "/web/chat.html")
				return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(text)}, nil
			}
			return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(config.MSG_PROCESSING)}, nil
		}
	})
	return res.(*message.Reply)
}

func (ex *OfficialAccount) processContent(ctx context.Context, msg *message.MixMessage, content *notion.Content, mr chan *message.Reply) {
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
	// set account info to content
	content.Account = accountInfo
	n := &notion.Notion{BearerToken: accountInfo.AccessToken, DatabaseID: accountInfo.DatabaseID}

	// 如果是不支持的类型，直接返回不支持的错误
	if _, ok := supportMsgTypeMap[msg.MsgType]; !ok {
		mr <- &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(config.MSG_UNSUPPOERT)}
		return
	}

	// 创建初始 Record
	res, pageID, err := n.CreateRecord(ctx, &notion.Content{
		Text:    "内容正在更新，请稍等",
		Account: accountInfo,
	})
	if err == nil {
		n.PageID = pageID
		go ex.updateNotionContent(ctx, msg, n, accountInfo, content)
	}
	mr <- &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(res)}
}

func updateLatestSchema(ctx context.Context, accountInfo *ent.Account, notionConfig *notion.Notion) {
	// 如果不是最新的 Scheam，更新 Schema
	if !accountInfo.IsLatestSchema {
		if _, err := notion.UpdateDatabaseProperties(ctx, notionConfig); err != nil {
			logger.SugaredLogger.Errorf("UpdateDatabaseProperties error: %s", err.Error())
		}
		_ = dao.UpdateIsLatestSchema(ctx, accountInfo.DatabaseID, true)
	}
}

func (ex *OfficialAccount) updateNotionContent(ctx context.Context, msg *message.MixMessage, n *notion.Notion, accountInfo *ent.Account, content *notion.Content) {
	ctx = context.WithValue(ctx, config.DATABASE_ID, accountInfo.DatabaseID)
	updateLatestSchema(ctx, accountInfo, n)
	content.Account = accountInfo
	content.NotionPageID = n.PageID
	content.Process(ctx)
	switch msg.MsgType {
	case message.MsgTypeText:
		// 保存文本信息到 Notion
		n.UpdateRecord(ctx, content)
	case message.MsgTypeImage, message.MsgTypeVideo, message.MsgTypeVoice:
		// 保存媒体信息到 Notion
		media := NewMedia(ex.officialAccount.GetContext())
		getMediaResp, err := media.getMedia(ctx, msg.MediaID, accountInfo.UserID)
		if err != nil {
			logger.SugaredLogger.Errorw("Get media from wechat error", "err", err)
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
