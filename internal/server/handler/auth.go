package handler

import (
	"context"
	"notionboy/internal/pkg/config"
	"time"

	model "notionboy/api/pb/model"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) GenrateToken(ctx context.Context, req *emptypb.Empty) (*model.GenrateTokenResponse, error) {
	tokenStr, err := s.AuthService.GenrateToken(ctx, "")
	if err != nil {
		return nil, err
	}

	return &model.GenrateTokenResponse{
		Token:  tokenStr,
		Type:   config.AUTH_HEADER_TOKEN_TYPE,
		Expiry: time.Now().Add(config.GetConfig().JWT.Expiration).Format(time.RFC3339),
	}, nil
}

func (s *Server) GenerateApiKey(ctx context.Context, req *emptypb.Empty) (*model.GenerateApiKeyResponse, error) {
	apiKey, err := s.AuthService.GenerateApiKey(ctx)
	if err != nil {
		return nil, err
	}
	return &model.GenerateApiKeyResponse{
		ApiKey: apiKey,
	}, nil
}

func (s *Server) DeleteApiKey(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	err := s.AuthService.DeleteApiKey(ctx)
	return &emptypb.Empty{}, err
}
