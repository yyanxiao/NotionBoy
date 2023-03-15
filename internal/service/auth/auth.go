package auth

import (
	"context"
	"fmt"

	"notionboy/db/ent"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/jwt"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/pkg/utils/cache"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var cacheClient = cache.DefaultClient()

type authServerImpl struct{}

type AuthServer interface {
	GenrateToken(ctx context.Context, userId, magicCode string) (string, error)
	GenerateApiKey(ctx context.Context) (string, error)
	DeleteApiKey(ctx context.Context) error
	GetUserIDFromContext(ctx context.Context) (uuid.UUID, error)
	GetAccountByApiKey(ctx context.Context, appiKey string) (*ent.Account, error)
	GetAccountByUserId(ctx context.Context, userId uuid.UUID) *ent.Account
	GetOAuthURL(ctx context.Context, provider string) (string, error)
	OAuthCallback(ctx context.Context, code string, state string) (string, error)
}

func NewAuthServer() AuthServer {
	return &authServerImpl{}
}

func (s *authServerImpl) GenrateToken(ctx context.Context, userId, magicCode string) (string, error) {
	var id uuid.UUID
	var err error
	logger.SugaredLogger.Debugw("GenrateToken", "userId", userId, "magicCode", magicCode)
	if magicCode != "" {
		val, ok := cacheClient.Get(fmt.Sprintf("%s:%s", config.MAGIC_CODE_CACHE_KEY, magicCode))
		if !ok {
			return "", status.Errorf(codes.Unauthenticated, "Invalid Magic Code")
		}
		acc := val.(*ent.Account)
		id = acc.UUID
	} else {
		if userId == "" {
			id, err = s.GetUserIDFromContext(ctx)
			if err != nil {
				return "", err
			}
		} else {
			id, err = uuid.Parse(userId)
			if err != nil {
				return "", status.Errorf(codes.Unauthenticated, "Invalid User")
			}
		}
		_, err = dao.QueryAccountByUUID(ctx, id)
		if err != nil {
			return "", err
		}
	}
	return jwt.GenerateToken(id.String())
}

func (s *authServerImpl) GenerateApiKey(ctx context.Context) (string, error) {
	userId, err := s.GetUserIDFromContext(ctx)
	if err != nil {
		return "", err
	}

	apiKey := uuid.New()
	if err := dao.UpdateAccountApiKey(ctx, userId, apiKey); err != nil {
		return "", err
	}
	return apiKey.String(), nil
}

func (s *authServerImpl) DeleteApiKey(ctx context.Context) error {
	userId, err := s.GetUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	return dao.UpdateAccountApiKey(ctx, userId, uuid.Nil)
}

func (s *authServerImpl) GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userId := ctx.Value(config.AUTH_USER_ID)

	if userId == nil {
		return uuid.Nil, status.Errorf(codes.Unauthenticated, "Invalid "+config.AUTH_HEADER_X_API_KEY)
	}
	id, err := uuid.Parse(userId.(string))
	if err != nil {
		return uuid.Nil, status.Errorf(codes.Unauthenticated, "Invalid User")
	}
	return id, nil
}

func (s *authServerImpl) GetAccountByApiKey(ctx context.Context, apiKeyStr string) (*ent.Account, error) {
	appiKey, err := uuid.Parse(apiKeyStr)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid "+config.AUTH_HEADER_X_API_KEY)
	}

	return dao.QueryAccountByApiKey(ctx, appiKey)
}

func (s *authServerImpl) GetAccountByUserId(ctx context.Context, userId uuid.UUID) *ent.Account {
	acc, err := dao.QueryAccountByUUID(ctx, userId)
	if err != nil {
		return nil
	}
	return acc
}
