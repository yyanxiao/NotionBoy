package handler

import (
	"context"
	"time"

	"notionboy/internal/pkg/config"

	model "notionboy/api/pb/model"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) GenrateToken(ctx context.Context, req *model.GenrateTokenRequest) (*model.GenrateTokenResponse, error) {
	tokenStr, err := s.AuthService.GenrateToken(ctx, "", req.GetMagicCode())
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

func (s *Server) OAuthCallback(ctx context.Context, req *model.OAuthCallbackRequest) (*model.GenrateTokenResponse, error) {
	tokenStr, err := s.AuthService.OAuthCallback(ctx, req.Code, req.State)
	if err != nil {
		return nil, err
	}
	return &model.GenrateTokenResponse{
		Token:  tokenStr,
		Type:   config.AUTH_HEADER_TOKEN_TYPE,
		Expiry: time.Now().Add(config.GetConfig().JWT.Expiration).Format(time.RFC3339),
	}, nil
}

func (s *Server) OAuthURL(ctx context.Context, req *model.OAuthURLRequest) (*model.OAuthURLResponse, error) {
	url, err := s.AuthService.GetOAuthURL(ctx, req.Provider)
	if err != nil {
		return nil, err
	}
	return &model.OAuthURLResponse{
		Url: url,
	}, nil
}
