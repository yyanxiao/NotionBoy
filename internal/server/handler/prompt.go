package handler

import (
	"context"

	"notionboy/api/pb/model"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) ListPrompts(ctx context.Context, req *emptypb.Empty) (*model.ListPromptsResponse, error) {
	return s.PromptService.ListPrompts(ctx)
}
