package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"notionboy/db/ent"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/service/auth"
	"notionboy/internal/service/completion"

	"github.com/sashabaranov/go-openai"
)

const OPENAI_HOST = "https://api.openai.com"

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

func Proxy(ctx context.Context, w http.ResponseWriter, r *http.Request, token string) {
	_, code, msg := checkAuthKey(ctx, token)
	if code != http.StatusOK {
		logger.SugaredLogger.Errorw("Validate api key error", "err", msg)
		http.Error(w, msg, code)
		return
	}

	url := OPENAI_HOST + r.URL.Path
	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		logger.SugaredLogger.Errorw("Error creating proxy request", "err", err)
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		return
	}
	for headerKey, headerValues := range r.Header {
		if headerKey == "Authorization" {
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.GetConfig().ChatGPT.ApiKey))
			continue
		}
		for _, headerValue := range headerValues {
			req.Header.Add(headerKey, headerValue)
		}
	}

	client := &http.Client{
		Timeout: 120 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.SugaredLogger.Errorw("Error proxying request", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// 将响应头复制到代理响应头中
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)

	// nolint:errcheck
	io.Copy(w, resp.Body)
}
