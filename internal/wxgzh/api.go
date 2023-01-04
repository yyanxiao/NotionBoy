package wxgzh

import (
	"fmt"
	"net/http"
	"notionboy/internal/chatgpt"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"
	"strings"

	notion "notionboy/internal/pkg/notion"

	wechat "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

// OfficialAccount 公众号操作样例
type OfficialAccount struct {
	wc              *wechat.Wechat
	officialAccount *officialaccount.OfficialAccount
	chatter         chatgpt.Chatter
}

// NewOfficialAccount
func NewOfficialAccount(wc *wechat.Wechat) *OfficialAccount {
	// init config
	wechatConfig := config.GetConfig().Wechat
	offCfg := &offConfig.Config{
		AppID:          wechatConfig.AppID,
		AppSecret:      wechatConfig.AppSecret,
		Token:          wechatConfig.Token,
		EncodingAESKey: wechatConfig.EncodingAESKey,
	}
	logger.SugaredLogger.Debug("AppID: ", offCfg.AppID)
	wc.SetCache(cache.NewMemory())
	officialAccount := wc.GetOfficialAccount(offCfg)
	return &OfficialAccount{
		wc:              wc,
		officialAccount: officialAccount,
		chatter:         chatgpt.DefaultApiClient(),
	}
}

func transformToNotionContent(msg *message.MixMessage) *notion.Content {
	content := notion.Content{
		Text: msg.Content,
	}
	return &content
}

// Serve 处理消息
func (ex *OfficialAccount) Serve(w http.ResponseWriter, r *http.Request) {
	// 传入request和responseWriter
	server := ex.officialAccount.GetServer(r, w)
	_, ok := server.GetQuery("signature")
	if !ok {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte("ok"))
		return
	}
	server.SkipValidate(false)
	// 设置接收消息的处理方法
	server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
		return ex.messageHandler(r.Context(), msg)
	})

	// 处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		if strings.Contains(err.Error(), "请求校验失败") {
			http.Error(w, fmt.Sprintf("400 Bad Request. err=%v", err), http.StatusBadRequest)
			return
		}
		http.Error(w, fmt.Sprintf("Serve Error, err=%v", err.Error()), http.StatusBadRequest)
		logger.SugaredLogger.Errorf("Serve Error, err=%v", err)
		return
	}
	// 发送回复的消息
	err = server.Send()
	if err != nil {
		logger.SugaredLogger.Errorf("Send Error, err=%v", err)
		return
	}
}

// GetAccessToken 获取ak
func (ex *OfficialAccount) GetAccessToken() (string, error) {
	ak, err := ex.officialAccount.GetAccessToken()
	if err != nil {
		logger.SugaredLogger.Errorf("get ak error, err=%v", err)
		return "", err
	}
	return ak, nil
}

// GetCallbackIP 获取微信callback IP地址
func (ex *OfficialAccount) GetCallbackIP() ([]string, error) {
	ipList, err := ex.officialAccount.GetBasic().GetCallbackIP()
	if err != nil {
		logger.SugaredLogger.Errorf("GetCallbackIP error, err=%v", err)
		return nil, err
	}
	return ipList, nil
}

// GetAPIDomainIP 获取微信callback IP地址
func (ex *OfficialAccount) GetAPIDomainIP() ([]string, error) {
	ipList, err := ex.officialAccount.GetBasic().GetAPIDomainIP()
	if err != nil {
		logger.SugaredLogger.Errorf("GetAPIDomainIP error, err=%v", err)
		return nil, err
	}
	return ipList, nil
}

// GetAPIDomainIP  清理接口调用次数
func (ex *OfficialAccount) ClearQuota() (string, error) {
	err := ex.officialAccount.GetBasic().ClearQuota()
	if err != nil {
		logger.SugaredLogger.Errorf("ClearQuota error, err=%v", err)
		return "", err
	}
	return "success", nil
}
