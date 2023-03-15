package auth

import (
	"context"
	"strings"

	"notionboy/internal/pkg/jwt"
	"notionboy/internal/pkg/logger"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getProvider(provider string) OauthProvider {
	var p OauthProvider
	switch provider {
	case "github":
		p = NewGithubProvider()
	}
	return p
}

func (s *authServerImpl) GetOAuthURL(ctx context.Context, provider string) (string, error) {
	p := getProvider(provider)
	if p == nil {
		return "", status.Errorf(codes.InvalidArgument, "invalid provider")
	}

	return p.GetOAuthURL(), nil
}

func (s *authServerImpl) OAuthCallback(ctx context.Context, code, state string) (string, error) {
	logger.SugaredLogger.Infow("OAuthCallback", "state", state, "code", code)
	states := strings.Split(state, ":")
	vender := states[0]
	provider := getProvider(vender)

	if provider == nil {
		return "", status.Errorf(codes.InvalidArgument, "invalid provider")
	}

	tok, err := provider.GetOAuthConf().Exchange(ctx, code)
	if err != nil {
		logger.SugaredLogger.Errorw("OAuthCallback failed", "error", err, "state", state, "provider", provider, "code", code, "token", tok)
		return "", err
	}

	acc, err := provider.QueryOrCreateNewUser(ctx, tok)
	if err != nil {
		return "", err
	}

	return jwt.GenerateToken(acc.UUID.String())
}
