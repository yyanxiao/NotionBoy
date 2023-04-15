package handler

import (
	"context"

	"notionboy/api/pb"
	"notionboy/internal/pkg/logger"

	model "notionboy/api/pb/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) CreateConversation(ctx context.Context, req *model.CreateConversationRequest) (*model.Conversation, error) {
	acc := getAccFromContext(ctx)
	if acc == nil {
		return nil, status.Errorf(codes.Unauthenticated, "Request unauthenticated")
	}

	dto, err := s.ConversationService.CreateConversation(ctx, acc, "", req.GetInstruction(), req.GetTitle())
	if err != nil {
		return nil, err
	}
	return dto.ToPB(), nil
}

func (s *Server) UpdateConversation(ctx context.Context, req *model.UpdateConversationRequest) (*model.Conversation, error) {
	acc := getAccFromContext(ctx)
	if acc == nil {
		return nil, status.Errorf(codes.Unauthenticated, "Request unauthenticated")
	}
	logger.SugaredLogger.Debugw("UpdateConversation", "req", req)
	dto, err := s.ConversationService.UpdateConversation(ctx, acc, req.Id, req.GetInstruction(), req.GetTitle())
	if err != nil {
		return nil, err
	}
	return dto.ToPB(), nil
}

func (s *Server) GetConversation(ctx context.Context, req *model.GetConversationRequest) (*model.Conversation, error) {
	acc := getAccFromContext(ctx)
	if acc == nil {
		return nil, status.Errorf(codes.Unauthenticated, "Request unauthenticated")
	}

	dto, err := s.ConversationService.GetConversation(ctx, acc, req.GetId())
	if err != nil {
		return nil, err
	}
	return dto.ToPB(), nil
}

func (s *Server) ListConversations(ctx context.Context, req *model.ListConversationsRequest) (*model.ListConversationsResponse, error) {
	acc := getAccFromContext(ctx)
	if acc == nil {
		return nil, status.Errorf(codes.Unauthenticated, "Request unauthenticated")
	}

	res, err := s.ConversationService.ListConversations(ctx, acc, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}
	var pbRes []*model.ConversationWithoutMessages
	for _, dto := range res {
		pbRes = append(pbRes, dto.ToPBWithoutMessages())
	}
	return &model.ListConversationsResponse{
		Conversations: pbRes,
	}, nil
}

func (s *Server) DeleteConversation(ctx context.Context, req *model.DeleteConversationRequest) (*emptypb.Empty, error) {
	acc := getAccFromContext(ctx)
	if acc == nil {
		return nil, status.Errorf(codes.Unauthenticated, "Request unauthenticated")
	}

	err := s.ConversationService.DeleteConversation(ctx, acc, req.GetId())
	return &emptypb.Empty{}, err
}

func (s *Server) CreateMessage(req *model.CreateMessageRequest, stream pb.Service_CreateMessageServer) error {
	ctx := stream.Context()
	acc := getAccFromContext(ctx)
	if acc == nil {
		return status.Errorf(codes.Unauthenticated, "Request unauthenticated")
	}

	err := s.ConversationService.CreateStreamConversationMessage(ctx, acc, stream, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) GetMessage(ctx context.Context, req *model.GetMessageRequest) (*model.Message, error) {
	acc := getAccFromContext(ctx)
	if acc == nil {
		return nil, status.Errorf(codes.Unauthenticated, "Request unauthenticated")
	}

	dto, err := s.ConversationService.GetConversationMessage(ctx, acc, req.GetConversationId(), req.GetId())
	if err != nil {
		return nil, err
	}
	return dto.ToPB(), nil
}

func (s *Server) ListMessages(ctx context.Context, req *model.ListMessagesRequest) (*model.ListMessagesResponse, error) {
	acc := getAccFromContext(ctx)
	if acc == nil {
		return nil, status.Errorf(codes.Unauthenticated, "Request unauthenticated")
	}

	res, err := s.ConversationService.ListConversationMessages(ctx, acc, req.GetConversationId(), int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}
	var pbRes []*model.Message
	for _, dto := range res {
		pbRes = append(pbRes, dto.ToPB())
	}
	return &model.ListMessagesResponse{
		Messages: pbRes,
	}, nil
}

func (s *Server) DeleteMessage(ctx context.Context, req *model.DeleteMessageRequest) (*emptypb.Empty, error) {
	acc := getAccFromContext(ctx)
	if acc == nil {
		return nil, status.Errorf(codes.Unauthenticated, "Request unauthenticated")
	}

	err := s.ConversationService.DeleteConversationMessage(ctx, acc, req.GetConversationId(), req.GetId())
	return &emptypb.Empty{}, err
}

func (s *Server) UpdateMessage(req *model.UpdateMessageRequest, stream pb.Service_UpdateMessageServer) error {
	ctx := stream.Context()
	acc := getAccFromContext(ctx)
	if acc == nil {
		return status.Errorf(codes.Unauthenticated, "Request unauthenticated")
	}

	err := s.ConversationService.UpdateStreamConversationMessage(ctx, acc, stream, req)
	if err != nil {
		return err
	}

	return nil
}
