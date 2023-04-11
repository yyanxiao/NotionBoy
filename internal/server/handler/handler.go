package handler

import (
	"notionboy/api/pb"
	"notionboy/internal/service/auth"
	"notionboy/internal/service/conversation"
	"notionboy/internal/service/order"
	"notionboy/internal/service/product"
	"notionboy/internal/service/prompt"
)

// Server is the gRPC server.
type Server struct {
	pb.UnimplementedServiceServer
	ConversationService conversation.ConversationService
	AuthService         auth.AuthServer
	OrderService        order.OrderService
	ProductService      product.ProductService
	PromptService       prompt.PromptService
}

func NewServer() *Server {
	return &Server{
		ConversationService: conversation.NewConversationService(),
		AuthService:         auth.NewAuthServer(),
		OrderService:        order.NewOrderService(),
		ProductService:      product.NewProductService(),
		PromptService:       prompt.NewPromptService(),
	}
}
