package notion

import (
	"context"
	"fmt"
	"strings"

	"notionboy/db/ent"
	"notionboy/db/ent/account"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/logger"

	"github.com/jomei/notionapi"
	"github.com/mitchellh/mapstructure"

	"golang.org/x/oauth2"
)

type OauthInterface interface {
	OAuthURL(state string) string
	OAuthProcess(state string) string
	OAuthCallback(ctx context.Context, code, state string) (string, error)
}

type oauthManager struct{}

func GetOauthManager() OauthInterface {
	return &oauthManager{}
}

func getOauthConf() *oauth2.Config {
	cfg := config.GetConfig().OAuth.Notion
	logger.SugaredLogger.Infof("oauthConf: %#v", cfg)

	return &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  cfg.AuthURL,
			TokenURL: cfg.AuthToken,
		},
	}
}

func (o *oauthManager) OAuthURL(state string) string {
	// url := "https://notionboy-test.theboys.tech/notion/oauth?state=" + state
	url := fmt.Sprintf("%s/notion/oauth?state=%s", config.GetConfig().Service.URL, state)
	logger.SugaredLogger.Debugf("Visit the OAuthURL: %v", url)
	return url
}

func (o *oauthManager) OAuthProcess(state string) string {
	oauthConf := getOauthConf()
	url := oauthConf.AuthCodeURL(state, oauth2.AccessTypeOffline)
	logger.SugaredLogger.Debugf("Visit the URL for the auth dialog: %v", url)
	return url
}

func (o *oauthManager) OAuthCallback(ctx context.Context, code, state string) (string, error) {
	states := strings.Split(state, ":")
	userType := states[0]
	userID := strings.Join(states[1:], "")
	oauthConf := getOauthConf()
	tok, err := oauthConf.Exchange(ctx, code)
	logger.SugaredLogger.Debugf("tok: %#v", tok)
	if err != nil {
		logger.SugaredLogger.Errorf("oauthConf.Exchange() failed with %v, code is %s", err, code)
		return "Get Oauth token failed", err
	}

	// oAuthInfo
	token := tok.AccessToken

	databaseID, err := bindNotion(ctx, token)
	if err != nil {
		logger.SugaredLogger.Errorf("GetDatabaseID() failed with %v", err)
		return "", err
	}

	acc := &ent.Account{
		UserID:         userID,
		UserType:       account.UserType(userType),
		AccessToken:    token,
		DatabaseID:     databaseID,
		IsLatestSchema: true,
	}
	// token extra: https://developers.notion.com/docs/authorization#step-4-notion-responds-with-an-access_token-and-some-additional-information
	if notionUser, err := parseUserInfo(tok.Extra("owner")); err == nil {
		acc.NotionUserID = notionUser.ID.String()
		acc.NotionUserEmail = notionUser.Person.Email

	}

	if err := dao.SaveAccount(ctx, acc); err != nil {
		logger.SugaredLogger.Errorw("Save account failed", "err", err, "account", acc)
		return "", err
	} else {
		return config.MSG_BIND_SUCCESS, nil
	}
}

func parseUserInfo(owner interface{}) (*notionapi.User, error) {
	user := owner.(map[string]interface{})["user"]
	var notionUser notionapi.User
	if err := mapstructure.Decode(user, &notionUser); err != nil {
		logger.SugaredLogger.Errorw("Get notion user info error", "err", err)
		return nil, err
	}
	return &notionUser, nil
}
