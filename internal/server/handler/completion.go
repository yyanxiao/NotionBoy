package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"notionboy/db/ent"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/service/auth"
	"notionboy/internal/service/completion"

	"github.com/sashabaranov/go-openai"
)

func Completions(ctx context.Context, w http.ResponseWriter, token string, req *openai.ChatCompletionRequest) {
	acc, code, msg := checkAuthKey(ctx, token)
	if code != http.StatusOK {
		logger.SugaredLogger.Errorw("Validate api key error", "err", msg)
		http.Error(w, msg, code)
		return
	}

	svc := completion.NewCompletionService()
	svc.Completions(ctx, w, acc, req)
}

func checkAuthKey(ctx context.Context, token string) (*ent.Account, int, string) {
	splits := strings.SplitN(token, " ", 2)
	if len(splits) < 2 || !strings.EqualFold(splits[0], "Bearer") {
		return nil, http.StatusUnauthorized, fmt.Sprintf("401 Unauthorized. %s", "invalid api key")
	}
	apiKey := splits[1]
	// check token
	acc, err := auth.NewAuthServer().GetAccountByApiKey(ctx, apiKey)
	if err != nil {
		logger.SugaredLogger.Errorw("auth failed", "token", token, "err", err)
		return nil, http.StatusForbidden, fmt.Sprintf("403 Forbidden. %s", err.Error())
	}
	if acc == nil {
		logger.SugaredLogger.Errorw("auth failed", "token", token, "err", err)
		return nil, http.StatusForbidden, "403 Forbidden. User not found"
	}
	return acc, http.StatusOK, "200 OK"
}
