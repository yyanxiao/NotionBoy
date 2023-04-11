package prompt

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"notionboy/api/pb/model"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/pkg/utils/cache"
)

var cacheClient = cache.DefaultClient()

const (
	enPromptURL = "https://ghproxy.com/https://raw.githubusercontent.com/f/awesome-chatgpt-prompts/main/prompts.csv"
	cnPromptURL = "https://ghproxy.com/https://raw.githubusercontent.com/PlexPt/awesome-chatgpt-prompts-zh/main/prompts-zh.json"
)

type PromptService interface {
	ListPrompts(context.Context) (*model.ListPromptsResponse, error)
}

type PromptServiceImpl struct{}

func NewPromptService() PromptService {
	return &PromptServiceImpl{}
}

func (s *PromptServiceImpl) ListPrompts(ctx context.Context) (*model.ListPromptsResponse, error) {
	enData, err := GetPromptsData(enPromptURL)
	if err != nil {
		return nil, err
	}
	cnData, err := GetPromptsData(cnPromptURL)
	if err != nil {
		return nil, err
	}

	data := append(enData, cnData...)

	return &model.ListPromptsResponse{
		Prompts: data,
	}, err
}

type Data struct {
	Act    string `json:"act" csv:"act"`
	Prompt string `json:"prompt" csv:"prompt"`
}

func GetPromptsData(url string) ([]*model.Prompt, error) {
	data, ok := cacheClient.Get(url)
	if ok {
		logger.SugaredLogger.Debugw("cache hit", "url", url)
		return data.([]*model.Prompt), nil
	}

	var prompts []*model.Prompt

	// Send a GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the content type of the response
	contentType := resp.Header.Get("Content-Type")
	urlStrs := strings.Split(url, ".")
	if urlStrs[len(urlStrs)-1] == "csv" {
		// If the content is CSV, read it using the CSV reader
		reader := csv.NewReader(resp.Body)
		records, err := reader.ReadAll()
		if err != nil {
			return nil, err
		}
		// Parse the CSV data into a slice of Data structs
		for _, row := range records[1:] { // Skip the header row
			if len(row) >= 2 {
				d := &model.Prompt{Act: row[0], Prompt: row[1]}
				prompts = append(prompts, d)
			}
		}
	} else if urlStrs[len(urlStrs)-1] == "json" {
		// If the content is JSON, read it using the JSON decoder
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(body, &prompts)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("unsupported file type: %s", contentType)
	}
	cacheClient.Set(url, prompts, 24*7*time.Hour)

	return prompts, nil
}
