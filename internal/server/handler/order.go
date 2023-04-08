package handler

import (
	"context"
	"net/http"

	"notionboy/internal/pkg/logger"
	"notionboy/internal/service/order"

	model "notionboy/api/pb/model"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (*model.Order, error) {
	acc := getAccFromContext(ctx)
	return s.OrderService.CreateOrder(ctx, acc, req)
}

func (s *Server) GetOrder(ctx context.Context, req *model.GetOrderRequest) (*model.Order, error) {
	acc := getAccFromContext(ctx)
	return s.OrderService.GetOrder(ctx, acc, req)
}

func (s *Server) ListOrders(ctx context.Context, req *model.ListOrdersRequest) (*model.ListOrdersResponse, error) {
	acc := getAccFromContext(ctx)
	logger.SugaredLogger.Debugw("ListOrders", "acc", acc)
	return s.OrderService.ListOrders(ctx, acc, req)
}

func (s *Server) DeleteOrder(ctx context.Context, req *model.DeleteOrderRequest) (*emptypb.Empty, error) {
	acc := getAccFromContext(ctx)
	return s.OrderService.DeleteOrder(ctx, acc, req)
}

// func (s *Server) UpdateOrder(ctx context.Context, req *model.UpdateOrderRequest) (*model.Order, error) {
// 	acc := getAccFromContext(ctx)
// 	return s.OrderService.UpdateOrder(ctx, acc, req)
// }

func (s *Server) PayOrder(ctx context.Context, req *model.PayOrderRequest) (*model.PayOrderResponse, error) {
	acc := getAccFromContext(ctx)
	return s.OrderService.PayOrder(ctx, acc, req)
}

func WechatPayCallback(ctx context.Context, r *http.Request) error {
	orderSvc := order.NewOrderService()
	return orderSvc.WechatPayCallBack(ctx, r)
}
