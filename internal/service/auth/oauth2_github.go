package auth

import (
	"context"
	"fmt"

	"notionboy/db/ent"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"

	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type GithubProvider struct {
	Name     string
	userType string
	State    string
}

func NewGithubProvider() OAuthProviderService {
	return &GithubProvider{
		Name:     PROVIDER_GITHUB,
		userType: PROVIDER_GITHUB,
		State:    config.GetConfig().OAuth.Github.State,
	}
}

func (o *GithubProvider) GetProviderName() string {
	return o.Name
}

func (o *GithubProvider) GetOAuthConf() *oauth2.Config {
	cfg := config.GetConfig().OAuth.Github
	return &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint:     github.Endpoint,
		Scopes:       []string{"user"},
	}
}

func (o *GithubProvider) GetOAuthURL() string {
	return o.GetOAuthConf().AuthCodeURL(fmt.Sprintf("%s:%s", o.userType, o.State))
}

func (o *GithubProvider) GetOAuthToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return o.GetOAuthConf().Exchange(ctx, code)
}

func (o *GithubProvider) QueryOrCreateNewUser(ctx context.Context, token *oauth2.Token) (*ent.Account, error) {
	user, err := getGithubUserInfo(token.AccessToken, token.TokenType)
	if err != nil {
		logger.SugaredLogger.Errorw("Failed to get github user info", "error", err)
		return nil, err
	}

	return queryOrCreateNewUser(ctx, user.Email, o.userType)
}

type GithubUser struct {
	Login       string `json:"login"`
	ID          int    `json:"id"`
	AvatarURL   string `json:"avatar_url"`
	URL         string `json:"url"`
	ReposURL    string `json:"repos_url"`
	Name        string `json:"name"`
	Company     string `json:"company"`
	Blog        string `json:"blog"`
	Location    string `json:"location"`
	Email       string `json:"email"`
	Bio         string `json:"bio"`
	PublicRepos int    `json:"public_repos"`
	PublicGists int    `json:"public_gists"`
	Followers   int    `json:"followers"`
	Following   int    `json:"following"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func getGithubUserInfo(token, tokenType string) (*GithubUser, error) {
	client := resty.New()
	var user GithubUser
	resp, err := client.R().
		SetResult(&user).
		SetAuthScheme(tokenType).
		SetAuthToken(token).
		Get("https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("Failed to get user info. Status Code: %d", resp.StatusCode())
	}

	return &user, nil
}
