package order

import (
	"context"
	"notionboy/api/pb/model"
	"notionboy/db/ent"

	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderServiceImpl struct{}

type OrderService interface {
	CreateOrder(ctx context.Context, acc *ent.Account, req *model.CreateOrderRequest) (*model.Order, error)
	GetOrder(ctx context.Context, acc *ent.Account, req *model.GetOrderRequest) (*model.Order, error)
	ListOrders(ctx context.Context, acc *ent.Account, req *model.ListOrdersRequest) (*model.ListOrdersResponse, error)
	DeleteOrder(ctx context.Context, acc *ent.Account, req *model.DeleteOrderRequest) (*emptypb.Empty, error)
	UpdateOrder(ctx context.Context, acc *ent.Account, req *model.UpdateOrderRequest) (*model.Order, error)
	PayOrder(ctx context.Context, acc *ent.Account, req *model.PayOrderRequest) (*model.PayOrderResponse, error)
}

func NewOrderService() OrderService {
	return &OrderServiceImpl{}
}
