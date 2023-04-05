package server

import (
	"encoding/json"
	"io"
	"net/http"

	"notionboy/internal/pkg/logger"
	"notionboy/internal/server/handler"

	"github.com/sashabaranov/go-openai"
)

func completions(w http.ResponseWriter, r *http.Request) {
	token := ""

	for k, v := range r.Header {
		if k == "Authorization" {
			token = v[0]
			break
		}
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logger.SugaredLogger.Errorw("Read request body failed", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req openai.ChatCompletionRequest

	err = json.Unmarshal(body, &req)
	if err != nil {
		logger.SugaredLogger.Errorw("Unmarshal request body failed", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	handler.Completions(r.Context(), w, token, &req)
}
