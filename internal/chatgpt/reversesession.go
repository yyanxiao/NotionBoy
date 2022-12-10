package chatgpt

import (
	"encoding/json"
	"net/http"
	"notionboy/internal/pkg/browser"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"
	"time"

	"github.com/go-resty/resty/v2"
)

var client *resty.Client

const (
	sessionURL         = "https://chat.openai.com/api/auth/session"
	loginUR            = "https://chat.openai.com/auth/login"
	cookieSessionToken = "__Secure-next-auth.session-token"
)

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

func (cli *reverseClient) login() {
	cfg := config.GetConfig().ChatGPT
	page := browser.New().MustConnect().MustPage(loginUR)
	page.MustElementR("button", "Log in").MustClick()

	page.MustElement("#username").MustInput(cfg.User)
	page.MustWaitLoad().MustElementR("button", "Continue").MustClick()
	page.MustElement("#password").MustInput(cfg.Pass)
	page.MustWaitLoad().MustElementR("button", "Continue").MustClick()
	time.Sleep(1 * time.Second)
	page.MustWaitLoad()
	cookies, _ := page.Cookies([]string{})
	for _, cookie := range cookies {
		if cookie.Name == cookieSessionToken {
			if cookie.Value != "" {
				logger.SugaredLogger.Info("Login to OpenAI use email success")
				config.GetConfig().ChatGPT.SessionToken = cookie.Value
			} else {
				logger.SugaredLogger.Warn("Login to OpenAI use email failed, did not get session token")
			}
			return
		}
	}
	logger.SugaredLogger.Warn("Login to OpenAI use email failed")
}

// RefreshSession use to keep session up to date
func (cli *reverseClient) refreshSession() {
	cli.setSessionTokenCookie()

	var resp *resty.Response
	var err error
	var newToken string

	resp, err = client.R().Get(sessionURL)
	if err != nil {
		logger.SugaredLogger.Errorw("refresh session for chatGPT error", "err", err)
		return
	}
	// if 401 return, token expired, need login to get a new token
	if resp.StatusCode() == http.StatusUnauthorized {
		cli.login()
		resp, err = client.R().Get(sessionURL)
		if err != nil {
			logger.SugaredLogger.Errorw("refresh session for chatGPT error", "err", err)
			return
		}
	}

	// if still don't get 200, set rate limit
	// and wait next time to retry
	if resp.StatusCode() != http.StatusOK {
		cli.setIsRateLimit(true)
		logger.SugaredLogger.Errorw("Refresh session for chatGPT error", "status", resp.Status())
		return
	}

	// update session token
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "__Secure-next-auth.session-token" {
			newToken = cookie.Value
			config.GetConfig().ChatGPT.SessionToken = cookie.Value
		}
	}

	// get auth token
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

	logger.SugaredLogger.Infow("refresh session success", "session_token", newToken)
	refreshHeaders()

	// if all pass, remove rate limit
	cli.setIsRateLimit(false)
}

func (cli *reverseClient) setSessionTokenCookie() {
	cfg := config.GetConfig().ChatGPT
	if cfg.SessionToken == "" && (cfg.User == "" || cfg.Pass == "") {
		logger.SugaredLogger.Fatal("Can not login to OpenAI, none of session token and username provided")
	}
	if cfg.SessionToken == "" {
		cli.login()
	}
	isSessionCookieExist := false
	for _, cookie := range client.Cookies {
		if cookie.Name == cookieSessionToken {
			cookie.Value = config.GetConfig().ChatGPT.SessionToken
			isSessionCookieExist = true
			break
		}
	}
	if !isSessionCookieExist {
		client.SetCookie(&http.Cookie{
			Name:  cookieSessionToken,
			Value: config.GetConfig().ChatGPT.SessionToken,
		})
	}
}
