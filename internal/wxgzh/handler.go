package wxgzh

import (
	"fmt"
	"notionboy/db/ent/account"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/db/dao"

	notion "notionboy/internal/pkg/notion"

	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	log "github.com/sirupsen/logrus"
)

func unBindingNotion(c *gin.Context, msg *message.MixMessage) *message.Reply {
	if err := dao.DeleteWxAccount(c, msg.GetOpenID()); err != nil {
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(config.MSG_UNBIND_FAILED + err.Error())}
	}
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(config.MSG_UNBIND_SUCCESS)}
}

func bindNotion(c *gin.Context, msg *message.MixMessage) *message.Reply {
	log.Warn("----- bindNotion ------")
	userID := msg.GetOpenID()
	userType := account.UserTypeWechat
	stage := fmt.Sprintf("%s:%s", userType, userID)
	oauthMgr := notion.GetOauthManager()
	url := oauthMgr.OAuthURL(stage)
	log.Info("OAuthURL: ", url)
	text := config.MSG_BINDING
	text += url
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(text)}
}

func helpInfo(c *gin.Context, msg *message.MixMessage) *message.Reply {
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(config.MSG_HELP)}
}
