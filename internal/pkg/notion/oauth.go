package notion

import (
	"context"
	"fmt"
	"notionboy/db/ent"
	"notionboy/db/ent/account"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/logger"
	"strings"

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
	logger.SugaredLogger.Infof("oauthConf: %#v", config.GetConfig().NotionOauth)
	return &oauth2.Config{
		ClientID:     config.GetConfig().NotionOauth.ClientID,
		ClientSecret: config.GetConfig().NotionOauth.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.GetConfig().NotionOauth.AuthURL,
			TokenURL: config.GetConfig().NotionOauth.AuthToken,
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

	if err := dao.SaveAccount(ctx, acc); err != nil {
		logger.SugaredLogger.Errorw("Save account failed", "err", err, "account", acc)
		return "", err
	} else {
		return config.MSG_BIND_SUCCESS, nil
	}
}
