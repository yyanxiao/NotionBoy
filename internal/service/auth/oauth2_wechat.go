package auth

import (
	"context"
	"fmt"
	"sync"

	"notionboy/db/ent"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"

	wechat "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/oauth"
	"golang.org/x/oauth2"
)

var once sync.Once

var wechatOA *officialaccount.OfficialAccount

type WeixinProvider struct {
	Name        string
	userType    string
	State       string
	RedirectUri string
	// officialAccount *officialaccount.OfficialAccount
	OAuthObj *oauth.Oauth
}

func getWechatOA() *officialaccount.OfficialAccount {
	if wechatOA == nil {
		once.Do(func() {
			wechatConfig := config.GetConfig().Wechat
			offCfg := &offConfig.Config{
				AppID:          wechatConfig.AppID,
				AppSecret:      wechatConfig.AppSecret,
				Token:          wechatConfig.Token,
				EncodingAESKey: wechatConfig.EncodingAESKey,
			}
			wc := wechat.NewWechat()
			wc.SetCache(cache.NewMemory())
			wechatOA = wc.GetOfficialAccount(offCfg)
		})
	}
	return wechatOA
}

func NewWeixinProvider() OAuthProviderService {
	officialAccount := getWechatOA()
	return &WeixinProvider{
		Name:        PROVIDER_WECHAT,
		userType:    PROVIDER_WECHAT,
		State:       config.GetConfig().OAuth.Wechat.State,
		RedirectUri: config.GetConfig().OAuth.Wechat.URLRedirect,
		OAuthObj:    oauth.NewOauth(officialAccount.GetContext()),
	}
}

func (o *WeixinProvider) GetProviderName() string {
	return o.Name
}

func (o *WeixinProvider) GetOAuthConf() *oauth2.Config {
	cfg := config.GetConfig().OAuth.Wechat
	return &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://open.weixin.qq.com/connect/oauth2/authorize",
			TokenURL: "https://api.weixin.qq.com/sns/oauth2/access_token",
		},
		Scopes: []string{"snsapi_userinfo"},
	}
}

func (o *WeixinProvider) GetOAuthURL() string {
	if uri, err := o.OAuthObj.GetRedirectURL(o.RedirectUri, o.GetOAuthConf().Scopes[0], fmt.Sprintf("%s:%s", o.userType, o.State)); err != nil {
		return ""
	} else {
		return uri
	}
}

func (o *WeixinProvider) GetOAuthToken(ctx context.Context, code string) (*oauth2.Token, error) {
	wtoken, err := o.OAuthObj.GetUserAccessToken(code)
	if err != nil {
		logger.SugaredLogger.Errorw("Failed to get Wechat access token", "error", err)
		return nil, err
	}
	token := &oauth2.Token{
		AccessToken:  wtoken.AccessToken,
		RefreshToken: wtoken.RefreshToken,
	}
	raw := make(map[string]interface{})
	raw["openid"] = wtoken.OpenID
	raw["unionid"] = wtoken.UnionID
	token = token.WithExtra(raw)
	logger.SugaredLogger.Debugw("GetOAuthToken", "token", token)
	return token, nil
}

func (o *WeixinProvider) QueryOrCreateNewUser(ctx context.Context, token *oauth2.Token) (*ent.Account, error) {
	logger.SugaredLogger.Debugw("QueryOrCreateNewUser", "token", token)
	openID := token.Extra("openid").(string)
	user, err := o.OAuthObj.GetUserInfo(token.AccessToken, openID, "zh_CN")
	logger.SugaredLogger.Debugw("QueryOrCreateNewUser", "user", user, "err", err)
	if err != nil {
		logger.SugaredLogger.Errorw("Failed to get github user info", "error", err)
		return nil, err
	}

	return queryOrCreateNewUser(ctx, user.OpenID, o.userType)
}
