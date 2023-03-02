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

// NewAuthFunc returns a new AuthFunc
// get key order:
// 1. get api key from header
// 2. get token from cookie
// 3. get Beare token from header
func NewAuthFunc() grpc_auth.AuthFunc {
	return func(c context.Context) (context.Context, error) {
		md := metautils.ExtractIncoming(c)
		// validate using api key
		val := md.Get(config.AUTH_HEADER_X_API_KEY)
		if val != "" {
			acc, err := auth.NewAuthServer().GetAccountByApiKey(c, val)
			if err != nil {
				return nil, err
			}
			// if account is exist, set account id to context
			newCtx := context.WithValue(c, config.ContextKeyUserId, acc.UUID)
			return newCtx, nil
		}

		// try to get token from cookie
		// if token is not exist, try to get token from Authorization header
		var err error
		token, ok := queryCookie(md, config.AUTH_HEADER_TOKEN)
		if !ok {
			// validate using jwt token
			token, err = grpc_auth.AuthFromMD(c, config.AUTH_HEADER_TOKEN_TYPE)
			if err != nil {
				return nil, err
			}
		}

		userId, err := jwt.ValidateToken(token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, err.Error())
		}

		uid, err := uuid.Parse(userId)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, err.Error())
		}

		acc := auth.NewAuthServer().GetAccountByUserId(c, uid)
		if acc == nil {
			return nil, status.Errorf(codes.Unauthenticated, "Invalid User")
		}

		newCtx := context.WithValue(c, config.ContextKeyUserId, userId)
		newCtx = context.WithValue(newCtx, config.ContextKeyUserAccount, acc)
		logger.SugaredLogger.Debugw("auth success", "userId", userId)
		return newCtx, nil
	}
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
