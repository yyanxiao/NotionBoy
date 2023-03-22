package auth

import (
	"context"
	"strings"

	"notionboy/db/ent"
	"notionboy/db/ent/account"
	"notionboy/internal/pkg/db/dao"

	"golang.org/x/oauth2"
)

type OAuthProviderService interface {
	GetProviderName() string
	GetOAuthConf() *oauth2.Config
	GetOAuthURL() string
	QueryOrCreateNewUser(ctx context.Context, token *oauth2.Token) (*ent.Account, error)
	GetOAuthToken(ctx context.Context, code string) (*oauth2.Token, error)
}

type OAuthProvider struct {
	Name        string
	userType    string
	State       string
	RedirectUri string
}

func (o *OAuthProvider) GetProviderName() string {
	return o.Name
}

func (o *OAuthProvider) GetUserType() string {
	return o.userType
}

func (o *OAuthProvider) GetOAuthConf() *oauth2.Config {
	return nil
}

func (o *OAuthProvider) GetOAuthURL() string {
	return ""
}

func (o *OAuthProvider) GetOAuthToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return nil, nil
}

type Provider struct {
	Name string
	URL  string
}

func queryOrCreateNewUser(ctx context.Context, userId, userType string) (*ent.Account, error) {
	var acc *ent.Account
	var err error
	// query account

	acc, err = dao.QueryAccount(ctx, account.UserType(userType), userId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			acc, err = dao.CreateAccount(ctx, &ent.Account{
				UserID:   userId,
				UserType: account.UserType(userType),
			})
			if err != nil {
				return nil, err
			}
			return acc, nil
		}
	}
	return acc, err
}
