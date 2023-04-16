package handler

import (
	"context"

	"notionboy/api/pb/model"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) ListPrompts(ctx context.Context, req *model.ListPromptsRequest) (*model.ListPromptsResponse, error) {
	acc := getAccFromContext(ctx)
	return s.PromptService.ListPrompts(ctx, acc, req)
}

func (s *Server) GetPrompt(ctx context.Context, req *model.GetPromptRequest) (*model.Prompt, error) {
	acc := getAccFromContext(ctx)
	return s.PromptService.GetPrompt(ctx, acc, req)
}

func (s *Server) CreatePrompt(ctx context.Context, req *model.CreatePromptRequest) (*model.Prompt, error) {
	acc := getAccFromContext(ctx)
	return s.PromptService.CreatePrompt(ctx, acc, req)
}

func (s *Server) UpdatePrompt(ctx context.Context, req *model.UpdatePromptRequest) (*model.Prompt, error) {
	acc := getAccFromContext(ctx)
	return s.PromptService.UpdatePrompt(ctx, acc, req)
}

func (s *Server) DeletePrompt(ctx context.Context, req *model.DeletePromptRequest) (*emptypb.Empty, error) {
	acc := getAccFromContext(ctx)
	return s.PromptService.DeletePrompt(ctx, acc, req)
}
