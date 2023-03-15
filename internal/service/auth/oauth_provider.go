package auth

import (
	"context"
	"strings"

	"notionboy/db/ent"
	"notionboy/db/ent/account"
	"notionboy/internal/pkg/db/dao"

	"golang.org/x/oauth2"
)

const PROVIDER_GITHUB = "github"

type OauthProvider interface {
	GetOAuthConf() *oauth2.Config
	GetOAuthURL() string
	QueryOrCreateNewUser(ctx context.Context, token *oauth2.Token) (*ent.Account, error)
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
