package auth

import (
	"context"
	"fmt"
	"notionboy/db/ent"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/jwt"
	"notionboy/internal/pkg/utils/cache"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/silenceper/wechat/v2/officialaccount/basic"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var cacheClient = cache.DefaultClient()

type authServerImpl struct{}

type AuthServer interface {
	GenrateToken(ctx context.Context, userId, magicCode, qrcode string) (string, error)
	GenerateApiKey(ctx context.Context) (string, error)
	DeleteApiKey(ctx context.Context) error
	GetUserIDFromContext(ctx context.Context) (uuid.UUID, error)
	GetAccountByApiKey(ctx context.Context, appiKey string) (*ent.Account, error)
	GetAccountByUserId(ctx context.Context, userId uuid.UUID) *ent.Account
	GetOAuthProviders(ctx context.Context) ([]Provider, error)
	OAuthCallback(ctx context.Context, code string, state string) (string, error)
	GenerateWechatQRCode(ctx context.Context) (string, string, error)
}

func NewAuthServer() AuthServer {
	return &authServerImpl{}
}

func (s *authServerImpl) GenrateToken(ctx context.Context, userId, magicCode, qrcode string) (string, error) {
	var id uuid.UUID
	var err error
	// logger.SugaredLogger.Debugw("GenrateToken", "userId", userId, "magicCode", magicCode, "qrcode", qrcode)
	if magicCode != "" {
		val, ok := cacheClient.Get(fmt.Sprintf("%s:%s", config.MAGIC_CODE_CACHE_KEY, magicCode))
		if !ok {
			return "", status.Errorf(codes.Unauthenticated, "Invalid Magic Code")
		}
		acc := val.(*ent.Account)
		id = acc.UUID
	} else if qrcode != "" {
		val, ok := cacheClient.Get(fmt.Sprintf("%s:%s", config.QRCODE_CACHE_KEY, qrcode))
		if !ok {
			return "", status.Errorf(codes.Unauthenticated, "Invalid QR Code")
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

func (s *authServerImpl) GetOAuthProviders(ctx context.Context) ([]Provider, error) {
	oauthProviderServices := []OAuthProviderService{
		NewGithubProvider(),
		NewWeixinProvider(),
	}
	ua := ctx.Value(config.ContextKeyUserAgent)

	var providers []Provider
	for _, oauthProviderService := range oauthProviderServices {
		// for wechat, only return wechat oauth url when user agent is wechat
		if oauthProviderService.GetProviderName() == PROVIDER_WECHAT {
			if ua != nil && strings.Contains(ua.(string), "MicroMessenger") {
				providers = append(providers, Provider{
					Name: oauthProviderService.GetProviderName(),
					URL:  oauthProviderService.GetOAuthURL(),
				})
			}
		} else {
			provider := Provider{
				Name: oauthProviderService.GetProviderName(),
				URL:  oauthProviderService.GetOAuthURL(),
			}
			providers = append(providers, provider)
		}
	}

	return providers, nil
}

func (s *authServerImpl) GenerateWechatQRCode(ctx context.Context) (string, string, error) {
	oa := getWechatOA()

	id := uuid.New().String()

	req := basic.NewTmpQrRequest(time.Duration(5)*time.Minute, id)
	// logger.SugaredLogger.Debugw("GenerateWechatQRCode", "req", req)
	b := basic.NewBasic(oa.GetContext())
	res, err := b.GetQRTicket(req)
	// logger.SugaredLogger.Debugw("GenerateWechatQRCode", "res", res, "err", err)
	if err != nil {
		return "", "", err
	}
	return fmt.Sprintf("https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=%s", res.Ticket), res.Ticket, nil
}
