package chatgpt

import (
	"encoding/json"
	"net/http"
	"notionboy/internal/pkg/browser"
	"notionboy/internal/pkg/logger"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	sessionURL         = "https://chat.openai.com/api/auth/session"
	loginUR            = "https://chat.openai.com/auth/login"
	cookieSessionToken = "__Secure-next-auth.session-token"
	userAgent          = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"
)

func (cli *reverseClient) login() {
	page := browser.New().MustConnect().MustPage(loginUR)
	page.MustElementR("button", "Log in").MustClick()
	page.MustElement("#username").MustInput(cli.Email)
	page.MustWaitLoad().MustElementR("button", "Continue").MustClick()
	page.MustElement("#password").MustInput(cli.Password)
	page.MustWaitLoad().MustElementR("button", "Continue").MustClick()
	time.Sleep(1 * time.Second)
	page.MustWaitLoad()
	cookies, _ := page.Cookies([]string{})
	for _, cookie := range cookies {
		if cookie.Name == cookieSessionToken {
			if cookie.Value != "" {
				logger.SugaredLogger.Info("Login to OpenAI use email success")
				cli.SessionToken = cookie.Value
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
	if cli.GetIsRateLimit() || cli.SessionToken == "" {
		cli.login()
	}

	var resp *resty.Response
	var err error

	resp, err = cli.client.R().
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", userAgent).
		SetCookie(&http.Cookie{
			Name:  cookieSessionToken,
			Value: cli.SessionToken,
		}).
		Get(sessionURL)
	if err != nil {
		logger.SugaredLogger.Errorw("refresh session for chatGPT error", "err", err)
		return
	}
	// if 401 return, token expired, need login to get a new token
	if resp.StatusCode() == http.StatusUnauthorized {
		cli.login()
		resp, err = cli.client.R().Get(sessionURL)
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
			cli.SessionToken = cookie.Value
			break
		}
	}

	// get auth token
	var data map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		logger.SugaredLogger.Errorw("Unmarshal refresh session for chatGPT error", "err", err)
		return
	}

	errMsg, ok := data["error"]
	if ok {
		logger.SugaredLogger.Errorw("refresh session for chatGPT error", "error", errMsg)
		return
	}

	accessToken, ok := data["accessToken"]
	if !ok {
		logger.SugaredLogger.Warn("Do not get token when refresh session for chatGPT")
		return
	}
	cli.authToken = accessToken.(string)

	logger.SugaredLogger.Infow("refresh session success", "session_token", cli.SessionToken)

	// if all pass, remove rate limit
	cli.setIsRateLimit(false)
}
