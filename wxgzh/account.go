package wxgzh

import (
	"notionboy/config"

	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	log "github.com/sirupsen/logrus"
)

//OfficialAccount 公众号
type OfficialAccount struct {
	wc              *wechat.Wechat
	officialAccount *officialaccount.OfficialAccount
}

//OfficialAccount 公众号操作实例
func NewOfficialAccount(wc *wechat.Wechat) *OfficialAccount {
	//init config
	wechatConfig := config.GetConfig().Wechat
	offCfg := &offConfig.Config{
		AppID:          wechatConfig.AppID,
		AppSecret:      wechatConfig.AppSecret,
		Token:          wechatConfig.Token,
		EncodingAESKey: wechatConfig.EncodingAESKey,
	}
	log.Debug("AppID: ", offCfg.AppID)
	wc.SetCache(cache.NewMemory())
	officialAccount := wc.GetOfficialAccount(offCfg)
	return &OfficialAccount{
		wc:              wc,
		officialAccount: officialAccount,
	}
}

//Serve 处理消息
func (ex *OfficialAccount) Serve(c *gin.Context) {
	// 传入request和responseWriter
	server := ex.officialAccount.GetServer(c.Request, c.Writer)
	server.SkipValidate(true)
	//设置接收消息的处理方法
	server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
		return messageHandler(c, msg)
	})

	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		log.Errorf("Serve Error, err=%v", err)
		return
	}
	//发送回复的消息
	err = server.Send()
	if err != nil {
		log.Errorf("Send Error, err=%v", err)
		return
	}
}
