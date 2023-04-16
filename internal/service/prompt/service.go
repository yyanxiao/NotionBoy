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
	"notionboy/db/ent"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/pkg/utils/cache"

	"github.com/google/uuid"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var cacheClient = cache.DefaultClient()

const (
	enPromptURL = "https://ghproxy.com/https://raw.githubusercontent.com/f/awesome-chatgpt-prompts/main/prompts.csv"
	cnPromptURL = "https://ghproxy.com/https://raw.githubusercontent.com/PlexPt/awesome-chatgpt-prompts-zh/main/prompts-zh.json"
)

func (s *PromptServiceImpl) ListPrompts(ctx context.Context, acc *ent.Account, req *model.ListPromptsRequest) (*model.ListPromptsResponse, error) {
	// list user custom prompts
	if req.IsCustom {
		prompts, err := dao.ListPrompts(ctx, acc.UUID)
		if err != nil {
			logger.SugaredLogger.Errorw("list prompts error", "err", err)
			return nil, err
		}
		var data []*model.Prompt
		for _, p := range prompts {
			dto := NewPromptDTO(p)
			data = append(data, dto.ToProto())
		}
		return &model.ListPromptsResponse{
			Prompts: data,
		}, err
	}
	// list default prompts
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

func (s *PromptServiceImpl) GetPrompt(ctx context.Context, acc *ent.Account, req *model.GetPromptRequest) (*model.Prompt, error) {
	pid, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(400, "invalid id, %s", err.Error())
	}
	p, err := dao.GetPrompt(ctx, pid, acc.UUID)
	if err != nil {
		logger.SugaredLogger.Errorw("get prompt error", "err", err, "id", req.Id)
		return nil, err
	}
	dto := NewPromptDTO(p)
	return dto.ToProto(), nil
}

func (s *PromptServiceImpl) CreatePrompt(ctx context.Context, acc *ent.Account, req *model.CreatePromptRequest) (*model.Prompt, error) {
	p, err := dao.CreatePrompt(ctx, acc.UUID, req.Act, req.Prompt)
	if err != nil {
		logger.SugaredLogger.Errorw("create prompt error", "err", err, "act", req.Act, "prompt", req.Prompt)
		return nil, err
	}
	dto := NewPromptDTO(p)
	return dto.ToProto(), nil
}

func (s *PromptServiceImpl) UpdatePrompt(ctx context.Context, acc *ent.Account, req *model.UpdatePromptRequest) (*model.Prompt, error) {
	pid, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(400, "invalid id, %s", err.Error())
	}
	p, err := dao.UpdatePrompt(ctx, pid, acc.UUID, req.Act, req.Prompt)
	if err != nil {
		logger.SugaredLogger.Errorw("update prompt error", "err", err, "id", pid, "act", req.Act, "prompt", req.Prompt, "user_id", acc.UUID)
		return nil, err
	}
	dto := NewPromptDTO(p)
	return dto.ToProto(), nil
}

func (s *PromptServiceImpl) DeletePrompt(ctx context.Context, acc *ent.Account, req *model.DeletePromptRequest) (*emptypb.Empty, error) {
	pid, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(400, "invalid id, %s", err.Error())
	}
	err = dao.DeletePrompt(ctx, pid, acc.UUID)
	return &emptypb.Empty{}, err
}

func GetPromptsData(url string) ([]*model.Prompt, error) {
	data, ok := cacheClient.Get(url)
	if ok {
		logger.SugaredLogger.Debugw("cache hit", "url", url)
		return data.([]*model.Prompt), nil
	}

	var prompts []*model.Prompt
	prompts = append(prompts, &model.Prompt{
		Id:       uuid.New().String(),
		Act:      "ChatGPT",
		Prompt:   "You are ChatGPT, a large language model trained by OpenAI. Follow the user's instructions carefully. Respond using markdown.",
		IsCustom: false,
	})

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
				d := &model.Prompt{
					Id:       uuid.New().String(),
					Act:      row[0],
					Prompt:   row[1],
					IsCustom: false,
				}
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
