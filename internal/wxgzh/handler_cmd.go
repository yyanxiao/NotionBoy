package wxgzh

import (
	"context"
	"fmt"
	"time"

	"notionboy/db/ent/account"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/pkg/utils/cache"
	"notionboy/internal/service/auth"

	notion "notionboy/internal/pkg/notion"

	"github.com/google/uuid"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

var cacheClient = cache.DefaultClient()

func unsubscribe(c context.Context, msg *message.MixMessage) {
	logger.SugaredLogger.Infow("unsubscribe", "openid", msg.GetOpenID())
	// if err := dao.DeleteAccount(c, account.UserTypeWechat, msg.GetOpenID()); err != nil {
	// 	logger.SugaredLogger.Errorw("delete account failed", "err", err)
	// }
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

func onSubscribe(c context.Context, msg *message.MixMessage) *message.Reply {
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(config.MSG_START)}
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

	token, err := svc.GenrateToken(ctx, acc.UUID.String(), "", "")
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

func magicCode(ctx context.Context, msg *message.MixMessage) *message.Reply {
	acc, err := dao.QueryAccountByWxUser(ctx, msg.GetOpenID())
	if err != nil {
		logger.SugaredLogger.Errorf("Query Account Error: %v", err)
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(config.MSG_ERROR_ACCOUNT_NOT_FOUND)}
	}
	code := uuid.New().String()
	cacheClient.Set(fmt.Sprintf("%s:%s", config.MAGIC_CODE_CACHE_KEY, code), acc, time.Duration(5)*time.Minute)
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(code)}
}

func scanQrcode(ctx context.Context, msg *message.MixMessage) *message.Reply {
	acc, err := dao.QueryAccountByWxUser(ctx, msg.GetOpenID())
	if err != nil {
		logger.SugaredLogger.Errorf("Query Account Error: %v", err)
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(config.MSG_ERROR_ACCOUNT_NOT_FOUND)}
	}

	id := msg.Ticket
	logger.SugaredLogger.Debugw("scan qrcode", "qrcode", id, "acc", acc, "msg", msg)
	cache.DefaultClient().Set(fmt.Sprintf("%s:%s", config.QRCODE_CACHE_KEY, id), acc, time.Duration(5)*time.Minute)
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText("登录成功")}
}

func whoAMI(ctx context.Context, msg *message.MixMessage) *message.Reply {
	myInfo, err := auth.WhoAmI(ctx, account.UserTypeWechat, msg.GetOpenID())
	if err != nil {
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(err.Error())}
	}
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(myInfo.String())}
}

func apiKey(ctx context.Context, msg *message.MixMessage) *message.Reply {
	key, err := auth.GenerateApiKey(ctx, account.UserTypeWechat, msg.GetOpenID())
	if err != nil {
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(err.Error())}
	}
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(key)}
}
