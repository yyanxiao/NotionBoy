package chatgpt

import (
	"context"
	"errors"
	"notionboy/db/ent"
	"notionboy/db/ent/quota"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/logger"
	"strings"
	"sync/atomic"

	gogpt "github.com/sashabaranov/go-gpt3"
)

type apiClient struct {
	*gogpt.Client
	isRateLimit atomic.Bool
}

func newApiClient(apiKey string) Chatter {
	client := &apiClient{
		Client: gogpt.NewClient(apiKey),
	}
	client.setIsRateLimit(false)
	return client
}

func (cli *apiClient) ChatWithHistory(ctx context.Context, acc *ent.Account, prompt string) (string, error) {
	isRateLimit, err := checkRateLimit(ctx, acc)
	if err != nil {
		return "", err
	}
	if isRateLimit {
		return "nil", errors.New(config.MSG_ERROR_QUOTA_LIMIT)
	}

	// get history
	h := getChatHistory(ctx, acc)
	h.Ctx = ctx
	h.Prompt = prompt
	newPrompt := h.BuildNewPrompt()

	// call api
	resp, err := cli.Chat(ctx, newPrompt)

	// process response
	h.Response = ProcessResponse(resp)
	h.MessageIdx++
	h.MessageID = resp.ID
	logger.SugaredLogger.Debugw("Response", "conversation_id", h.MessageID, "error", nil, "message", h.Response, "usage", resp.Usage)
	if err == nil {
		// save history to cache and db
		setChatHistory(h)
		err = h.SaveHistory()
		if err != nil {
			logger.SugaredLogger.Errorf("Save history error: %v", err)
			return h.Response, err
		}
		err = dao.IncrDailyQuota(ctx, acc.ID, quota.CategoryChatgpt)
		if err != nil {
			logger.SugaredLogger.Errorf("Save history error: %v", err)
			return h.Response, err
		}
	}
	return h.Response, err
}

func checkRateLimit(ctx context.Context, acc *ent.Account) (bool, error) {
	qt, err := dao.QueryQuota(ctx, acc.ID, quota.CategoryChatgpt)
	if err != nil {
		logger.SugaredLogger.Errorf("Query Quota Error: %v", err)
		return false, err
	}
	if qt.DailyUsed >= qt.Daily {
		logger.SugaredLogger.Debugw("Hit rate limit", "account", acc.ID, "daily_used", qt.DailyUsed, "daily", qt.Daily)
		return true, nil
	}
	logger.SugaredLogger.Debugw("Not hit rate limit", "account", acc.ID, "daily_used", qt.DailyUsed, "daily", qt.Daily, "category", qt.Category)

	return false, nil
}

func (cli *apiClient) Chat(ctx context.Context, prompt string) (*gogpt.CompletionResponse, error) {
	if cli.GetIsRateLimit() {
		return nil, errors.New("hit rate limit, please increase your quote")
	}
	logger.SugaredLogger.Debugw("Get prompt message for api client", "prompt", prompt)
	model := gogpt.GPT3TextDavinci003
	if config.GetConfig().ChatGPT.Model != "" {
		model = config.GetConfig().ChatGPT.Model
	}
	req := gogpt.CompletionRequest{
		Model:       model,
		MaxTokens:   2048,
		Prompt:      prompt,
		Temperature: 0.5,
	}

	respChan := make(chan *gogpt.CompletionResponse)
	errChan := make(chan error)

	chat := func() {
		resp, err := cli.CreateCompletion(ctx, req)
		if err != nil {
			errChan <- err
		} else {
			respChan <- &resp
		}
	}
	var err error
	for i := 0; i < 3; i++ {
		go chat()
		select {
		case resp := <-respChan:
			return resp, nil
			// msgId := resp.ID
			// sb := strings.Builder{}

			// for _, item := range resp.Choices {
			// 	sb.WriteString(item.Text)
			// 	sb.WriteString("\n")
			// }
			// logger.SugaredLogger.Debugw("Response", "conversation_id", msgId, "error", nil, "message", sb.String(), "usage", resp.Usage)
			// return processChatResponse(sb.String()), nil
		case err = <-errChan:
			logger.SugaredLogger.Warnw("Get response from chatGPT error", "retry_times", i+1, "err", err)
		}
	}

	return nil, err
}

func (cli *apiClient) ResetHistory(acc *ent.Account) {
	resetChatHistory(acc)
}

func (cli *apiClient) GetIsRateLimit() bool {
	return cli.isRateLimit.Load()
}

func (cli *apiClient) setIsRateLimit(flag bool) {
	cli.isRateLimit.Store(flag)
}

// ProcessResponse get response string from gogpt.CompletionResponse
func ProcessResponse(resp *gogpt.CompletionResponse) string {
	sb := strings.Builder{}
	for _, item := range resp.Choices {
		sb.WriteString(item.Text)
		sb.WriteString("\n")
	}
	s := sb.String()
	s = strings.Replace(s, "<|im_start|>", "", -1)
	s = strings.Replace(s, "<|im_end|>", "", -1)
	s = strings.TrimSpace(s)
	return s
}
