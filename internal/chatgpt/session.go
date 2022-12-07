package chatgpt

import (
	"encoding/json"
	"net/http"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"

	"github.com/go-resty/resty/v2"
)

var client *resty.Client

const authURL = "https://chat.openai.com/api/auth/session"

func init() {
	client = resty.New()
	refreshHeaders()
}

func refreshHeaders() {
	cfg := config.GetConfig().ChatGPT
	client.SetHeaders(map[string]string{
		"Accept":        "application/json",
		"Authorization": "Bearer " + cfg.Authorization,
		"Content-Type":  "application/json",
		"User-Agent":    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
	})
}

// RefreshSession use to keep session up to date
func refreshSession() {
	cfg := config.GetConfig().ChatGPT
	if cfg.SessionToken == "" {
		logger.SugaredLogger.Fatal("Can't find sessionToken")
	}
	client.SetCookie(&http.Cookie{
		Name:  "__Secure-next-auth.session-token",
		Value: cfg.SessionToken,
	})

	resp, err := client.R().Get(authURL)
	if err != nil {
		logger.SugaredLogger.Errorw("refresh session for chatGPT error", "err", err)
		return
	}
	if resp.StatusCode() != http.StatusOK {
		logger.SugaredLogger.Errorw("Refresh session for chatGPT error", "status", resp.Status())
		return
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "__Secure-next-auth.session-token" {
			config.GetConfig().ChatGPT.SessionToken = cookie.Value
		}
	}
	var data map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		logger.SugaredLogger.Errorw("Unmarshal refresh session for chatGPT error", "err", err)
		return
	}
	accessToken, ok := data["accessToken"]
	if !ok {
		logger.SugaredLogger.Warn("Do not get token when refresh session for chatGPT")
		return
	}
	config.GetConfig().ChatGPT.Authorization = accessToken.(string)

	logger.SugaredLogger.Debugf("%#v", cfg)
	refreshHeaders()
}
