package handler

import (
	"notionboy/api/pb"
	"notionboy/internal/service/auth"
	"notionboy/internal/service/conversation"
)

// Server is the gRPC server.
type Server struct {
	pb.UnimplementedServiceServer
	ConversationService conversation.ConversationService
	AuthService         auth.AuthServer
}

func NewServer() *Server {
	return &Server{
		ConversationService: conversation.NewConversationService(),
		AuthService:         auth.NewAuthServer(),
	}
}
