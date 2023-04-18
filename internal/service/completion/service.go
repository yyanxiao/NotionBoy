package completion

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"notionboy/db/ent"
	"notionboy/internal/pkg/db"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/logger"

	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
)

func (s *CompletionServiceImpl) Completions(ctx context.Context, w http.ResponseWriter, acc *ent.Account, req *openai.ChatCompletionRequest) {
	isRateLimited, err := checkQuota(ctx, acc, req)
	if err != nil {
		logger.SugaredLogger.Errorw("Proxy Completions error", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		// nolint:errcheck
		w.Write([]byte(err.Error()))
		return
	}
	if isRateLimited {
		w.WriteHeader(http.StatusTooManyRequests)
		// nolint:errcheck
		w.Write([]byte("额度已用完，请点击公众号菜单栏服务中的 VIP 进行充值"))
		return
	}

	if req.Stream {
		s.chatStream(ctx, w, acc, req)
	} else {
		s.chat(ctx, w, acc, req)
	}
}

func (s *CompletionServiceImpl) chat(ctx context.Context, w http.ResponseWriter, acc *ent.Account, req *openai.ChatCompletionRequest) {
	resp, err := s.Client.CreateChatCompletion(ctx, *req)
	if err != nil {
		logger.SugaredLogger.Errorw("CreateChatCompletion error", "err", err, "resp", resp)
		w.WriteHeader(http.StatusInternalServerError)
		// nolint:errcheck
		w.Write([]byte(err.Error()))
		return
	}
	jsonData, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	updateTokenUsage(ctx, acc, int64(resp.Usage.TotalTokens))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// nolint:errcheck
	w.Write(jsonData)
}

func (s *CompletionServiceImpl) chatStream(ctx context.Context, w http.ResponseWriter, acc *ent.Account, req *openai.ChatCompletionRequest) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	stream, err := s.Client.CreateChatCompletionStream(ctx, *req)
	if err != nil {
		logger.SugaredLogger.Errorw("CreateChatCompletion error", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		// nolint:errcheck
		w.Write([]byte(err.Error()))
		return
	}
	defer stream.Close()
	sb := strings.Builder{}

	for {
		response, err := stream.Recv()
		// logger.SugaredLogger.Debugw("stream.Recv", "response", response, "err", err)

		if errors.Is(err, io.EOF) {
			// stream is done
			streamDataToClient(w, []byte("[DONE]"))
			tokenUsage := calcuteTokenUsage(req, sb.String())
			updateTokenUsage(ctx, acc, tokenUsage)
			return
		}
		// save the last response to calcute token usage
		sb.WriteString(response.Choices[0].Delta.Content)

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return
		}
		message, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			// nolint:errcheck
			w.Write([]byte(err.Error()))
			return
		}
		streamDataToClient(w, message)
	}
}

func streamDataToClient(w http.ResponseWriter, message []byte) {
	fmt.Fprintf(w, "data: %s\n\n", message)
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	} else {
		logger.SugaredLogger.Errorw("Streaming unsupported!")
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
}

func calcuteTokenUsage(req *openai.ChatCompletionRequest, respStr string) int64 {
	totalTokens := 0
	promptTokens := 0
	completionTokens := 0
	tk, _ := tiktoken.GetEncoding("cl100k_base")

	// token for req
	for _, msg := range req.Messages {
		promptTokens += len(tk.Encode(msg.Content, nil, nil))
	}

	// token for resp
	completionTokens += len(tk.Encode(respStr, nil, nil))

	usage := &openai.Usage{
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
		TotalTokens:      promptTokens + completionTokens,
	}

	switch req.Model {
	case openai.GPT3Dot5Turbo, openai.GPT3Dot5Turbo0301:
		totalTokens = usage.TotalTokens
	case openai.GPT4, openai.GPT40314:
		totalTokens = usage.PromptTokens*15 + usage.CompletionTokens*30
	case openai.GPT432K, openai.GPT432K0314:
		totalTokens = usage.PromptTokens*30 + usage.CompletionTokens*60
	}
	return int64(totalTokens)
}

func updateTokenUsage(ctx context.Context, acc *ent.Account, usage int64) {
	// logger.SugaredLogger.Debugw("updateTokenUsage", "acc", acc, "usage", usage)
	if err := dao.IncrUsedTokenQuota(db.GetClient(), ctx, acc.ID, usage); err != nil {
		logger.SugaredLogger.Errorw("IncrUsedTokenQuota error", "err", err)
	}
}

// checkQuota check if the account is rate limited
func checkQuota(ctx context.Context, acc *ent.Account, req *openai.ChatCompletionRequest) (bool, error) {
	quota, err := dao.QueryQuota(ctx, acc.ID)
	if err != nil {
		return true, err
	}
	if quota.TokenUsed >= quota.Token {
		return true, nil
	}

	tokens := calcuteTokenUsage(req, "")
	switch req.Model {
	case openai.GPT4, openai.GPT40314:
		if ((tokens + 1000) * 15) > (quota.Token - quota.TokenUsed) {
			return true, nil
		}

	case openai.GPT432K, openai.GPT432K0314:
		if ((tokens + 1000) * 30) > (quota.Token - quota.TokenUsed) {
			return true, nil
		}
	default:
		if ((tokens + 1000) * 1) > (quota.Token - quota.TokenUsed) {
			return true, nil
		}
	}
	return false, nil
}
