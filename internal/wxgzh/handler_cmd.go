package wxgzh

import (
	"context"
	"fmt"
	"notionboy/db/ent/account"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/service/auth"

	notion "notionboy/internal/pkg/notion"

	"github.com/silenceper/wechat/v2/officialaccount/message"
)

func unsubscribe(c context.Context, msg *message.MixMessage) {
	if err := dao.DeleteAccount(c, account.UserTypeWechat, msg.GetOpenID()); err != nil {
		logger.SugaredLogger.Errorw("delete account failed", "err", err)
	}
}

func unBindingNotion(c context.Context, msg *message.MixMessage) *message.Reply {
	if err := dao.ClearNotionAuthInfo(c, account.UserTypeWechat, msg.GetOpenID()); err != nil {
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(config.MSG_UNBIND_FAILED + err.Error())}
	}
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(config.MSG_UNBIND_SUCCESS)}
}

func bindNotion(c context.Context, msg *message.MixMessage) *message.Reply {
	logger.SugaredLogger.Info("----- bindNotion ------")
	userID := msg.GetOpenID()
	userType := account.UserTypeWechat
	stage := fmt.Sprintf("%s:%s", userType, userID)
	oauthMgr := notion.GetOauthManager()
	url := oauthMgr.OAuthURL(stage)
	logger.SugaredLogger.Info("OAuthURL: ", url)
	text := config.MSG_BINDING
	text += url
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(text)}
}

func helpInfo(c context.Context, msg *message.MixMessage) *message.Reply {
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(config.MSG_HELP)}
}

func sosInfo(c context.Context, msg *message.MixMessage) *message.Reply {
	return &message.Reply{
		MsgType: message.MsgTypeText,
		MsgData: message.NewText(fmt.Sprintf("欢迎添加作者微信，请搜索🔍:  %s", config.GetConfig().Wechat.AuthorID)),
	}
}

func webui(ctx context.Context, msg *message.MixMessage) *message.Reply {
	acc, err := dao.QueryAccountByWxUser(ctx, msg.GetOpenID())
	if err != nil {
		return &message.Reply{
			MsgType: message.MsgTypeText,
			MsgData: message.NewText(fmt.Sprintf("查询账户信息失败: %s", err.Error())),
		}
	}

	svc := auth.NewAuthServer()

	token, err := svc.GenrateToken(ctx, acc.UUID.String())
	if err != nil {
		return &message.Reply{
			MsgType: message.MsgTypeText,
			MsgData: message.NewText(fmt.Sprintf("生成 Token 失败: %s", err.Error())),
		}
	}

	webui := fmt.Sprintf("%s/web?token=%s", config.GetConfig().Service.URL, token)

	return &message.Reply{
		MsgType: message.MsgTypeText,
		MsgData: message.NewText(fmt.Sprintf("欢迎访问 NotionBoy 的 WebUI: %s", webui)),
	}
}
