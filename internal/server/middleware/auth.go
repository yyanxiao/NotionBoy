package middleware

import (
	"context"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/jwt"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/service/auth"
	"strings"

	"github.com/google/uuid"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var skipAuthPaths = []string{
	"/v1/auth/token",
	"/v1/auth/callback",
	"/v1/auth/providers",
	"/v1/auth/wechat/qrcode",
	"/v1/products",
}

// NewAuthFunc returns a new AuthFunc
// get key order:
// 1. get api key from header
// 2. get token from cookie
// 3. get Beare token from header
// 4. skip auth for some paths
func NewAuthFunc() grpc_auth.AuthFunc {
	return func(c context.Context) (context.Context, error) {
		md := metautils.ExtractIncoming(c)
		// logger.SugaredLogger.Debugw("auth", "md", md)

		ua := md.Get(config.ContextKeyUserAgent.String())
		if ua != "" {
			c = context.WithValue(c, config.ContextKeyUserAgent, ua)
		}

		// validate using api key
		apiKey := md.Get(config.AUTH_HEADER_X_API_KEY)
		if apiKey != "" {
			// logger.SugaredLogger.Debugw("auth by api key", "apiKey", apiKey)
			acc, err := auth.NewAuthServer().GetAccountByApiKey(c, apiKey)
			logger.SugaredLogger.Debugw("auth by api key", "apiKey", apiKey, "account", acc, "err", err)
			if err != nil {
				return nil, err
			}
			if acc == nil {
				return nil, status.Errorf(codes.Unauthenticated, "Invalid User")
			}
			// if account is exist, set account id to context
			newCtx := context.WithValue(c, config.ContextKeyUserId, acc.UUID)
			newCtx = context.WithValue(newCtx, config.ContextKeyUserAccount, acc)
			return newCtx, nil
		}

		// validate using token
		cookieToken, hasCookieToken := queryCookie(md, config.AUTH_HEADER_TOKEN)

		bearerToken, err := authFromMD(c, config.AUTH_HEADER_TOKEN_TYPE)
		if err != nil {
			return nil, err
		}
		if hasCookieToken || bearerToken != "" {
			newCtx, err := validateByToken(c, cookieToken, bearerToken)
			if err == nil {
				return newCtx, nil
			}
		}

		// skip auth for some paths
		path := md.Get("path")
		if path != "" {
			for _, skipPath := range skipAuthPaths {
				if strings.HasPrefix(path, skipPath) {
					return c, nil
				}
			}
		}

		// if no api key, cookie token, bearer token and not skip auth path return error
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}
}

func validateByToken(ctx context.Context, cookieToken, bearerToken string) (context.Context, error) {
	var token string
	if cookieToken != "" {
		token = cookieToken
	} else {
		token = bearerToken
	}
	if token == "" {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid Token")
	}
	// logger.SugaredLogger.Debugw("auth by token", "token", token)

	userId, err := jwt.ValidateToken(token)
	if err != nil {
		logger.SugaredLogger.Errorw("auth failed", "token", token, "err", err)
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	uid, err := uuid.Parse(userId)
	if err != nil {
		logger.SugaredLogger.Errorw("auth failed with user id", "token", token, "err", err, "userId", userId)
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	acc := auth.NewAuthServer().GetAccountByUserId(ctx, uid)
	if acc == nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid User")
	}

	newCtx := context.WithValue(ctx, config.ContextKeyUserId, userId)
	newCtx = context.WithValue(newCtx, config.ContextKeyUserAccount, acc)
	// logger.SugaredLogger.Debugw("auth success", "userId", userId)
	return newCtx, nil
}

func queryCookie(md metautils.NiceMD, name string) (string, bool) {
	ck := md.Get(config.AUTH_HEADER_COOKIE)
	if ck == "" {
		return "", false
	}
	cookies := strings.Split(ck, ";")
	for _, c := range cookies {
		cookie := strings.Split(c, "=")
		if len(cookie) == 2 {
			key := strings.TrimSpace(cookie[0])
			val := strings.TrimSpace(cookie[1])
			if key == name {
				return val, true
			}
		}
	}
	return "", false
}

func authFromMD(ctx context.Context, expectedScheme string) (string, error) {
	val := metautils.ExtractIncoming(ctx).Get("authorization")
	if val == "" {
		return "", nil
	}
	splits := strings.SplitN(val, " ", 2)
	if len(splits) < 2 {
		return "", status.Errorf(codes.Unauthenticated, "Bad authorization string")
	}
	if !strings.EqualFold(splits[0], expectedScheme) {
		return "", status.Errorf(codes.Unauthenticated, "Request unauthenticated with "+expectedScheme)
	}
	return splits[1], nil
}
