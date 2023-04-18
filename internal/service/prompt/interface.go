package prompt

import (
	"context"

	"notionboy/api/pb/model"
	"notionboy/db/ent"

	"google.golang.org/protobuf/types/known/emptypb"
)

type PromptService interface {
	ListPrompts(context.Context, *ent.Account, *model.ListPromptsRequest) (*model.ListPromptsResponse, error)
	GetPrompt(context.Context, *ent.Account, *model.GetPromptRequest) (*model.Prompt, error)
	CreatePrompt(context.Context, *ent.Account, *model.CreatePromptRequest) (*model.Prompt, error)
	UpdatePrompt(context.Context, *ent.Account, *model.UpdatePromptRequest) (*model.Prompt, error)
	DeletePrompt(context.Context, *ent.Account, *model.DeletePromptRequest) (*emptypb.Empty, error)
}

type PromptServiceImpl struct{}

func NewPromptService() PromptService {
	return &PromptServiceImpl{}
}
