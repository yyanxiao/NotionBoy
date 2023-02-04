package chatgpt

import (
	"context"
	"errors"
	"notionboy/db/ent"
	"notionboy/internal/pkg/config"
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
	newPrompt := buildPrompt(ctx, acc, prompt)
	resp, err := cli.Chat(ctx, newPrompt)
	if err == nil {
		setChatHistory(ctx, acc, prompt, resp)
	}
	return resp, err
}

func (cli *apiClient) Chat(ctx context.Context, prompt string) (string, error) {
	if cli.GetIsRateLimit() {
		return "", errors.New("hit rate limit, please increase your quote")
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
		Temperature: 0.9,
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
			msgId := resp.ID
			sb := strings.Builder{}

			for _, item := range resp.Choices {
				sb.WriteString(item.Text)
				sb.WriteString("\n")
			}
			logger.SugaredLogger.Debugw("Response", "conversation_id", msgId, "error", nil, "message", sb.String(), "usage", resp.Usage)
			return processChatResponse(sb.String()), nil
		case err = <-errChan:
			logger.SugaredLogger.Warnw("Get response from chatGPT error", "retry_times", i+1, "err", err)
		}
	}

	return "", err
}

func (cli *apiClient) GetIsRateLimit() bool {
	return cli.isRateLimit.Load()
}

func (cli *apiClient) setIsRateLimit(flag bool) {
	cli.isRateLimit.Store(flag)
}

func processChatResponse(resp string) string {
	resp = strings.Replace(resp, "<|im_start|>", "", -1)
	resp = strings.Replace(resp, "<|im_end|>", "", -1)
	resp = strings.TrimSpace(resp)
	return resp
}
