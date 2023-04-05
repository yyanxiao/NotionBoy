package handler

import (
	"context"
	"time"

	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"

	model "notionboy/api/pb/model"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) GenrateToken(ctx context.Context, req *model.GenrateTokenRequest) (*model.GenrateTokenResponse, error) {
	tokenStr, err := s.AuthService.GenrateToken(ctx, "", req.GetMagicCode(), req.Qrcode)
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

func (s *Server) OAuthProviders(ctx context.Context, req *model.OAuthURLRequest) (*model.OAuthURLResponse, error) {
	providers, err := s.AuthService.GetOAuthProviders(ctx)
	if err != nil {
		return nil, err
	}

	resp := make([]*model.OAuthProvider, 0, len(providers))
	for _, p := range providers {
		resp = append(resp, &model.OAuthProvider{
			Name: p.Name,
			Url:  p.URL,
		})
	}
	return &model.OAuthURLResponse{
		Providers: resp,
	}, nil
}

func (s *Server) GenerateWechatQRCode(ctx context.Context, req *emptypb.Empty) (*model.GenerateWechatQRCodeResponse, error) {
	url, code, err := s.AuthService.GenerateWechatQRCode(ctx)
	logger.SugaredLogger.Debugw("GenerateWechatQRCode", "url", url, "code", code, "err", err)
	if err != nil {
		return nil, err
	}
	return &model.GenerateWechatQRCodeResponse{
		Url:    url,
		Qrcode: code,
	}, nil
}
