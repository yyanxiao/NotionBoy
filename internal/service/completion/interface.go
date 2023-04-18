package completion

import (
	"context"
	"net/http"
	"sync"

	"notionboy/db/ent"
	"notionboy/internal/pkg/config"

	"github.com/sashabaranov/go-openai"
)

var (
	client *openai.Client
	once   sync.Once
)

type CompletionService interface {
	Completions(ctx context.Context, w http.ResponseWriter, acc *ent.Account, req *openai.ChatCompletionRequest)
}

type CompletionServiceImpl struct {
	Client *openai.Client
}

func NewCompletionService() CompletionService {
	if client == nil {
		once.Do(func() {
			client = openai.NewClient(config.GetConfig().ChatGPT.ApiKey)
		})
	}

	return &CompletionServiceImpl{Client: client}
}
