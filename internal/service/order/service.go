package order

import (
	"context"
	"net/http"
	"strings"

	"notionboy/api/pb/model"
	"notionboy/db/ent"
	"notionboy/db/ent/order"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/db"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/service/order/pay"

	"github.com/google/uuid"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CreateOrder creates a new order
func (s *OrderServiceImpl) CreateOrder(ctx context.Context, acc *ent.Account, req *model.CreateOrderRequest) (*model.Order, error) {
	pid, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Errorf(400, "invalid id, %s", err.Error())
	}
	p, err := dao.GetProduct(ctx, pid)
	if err != nil {
		return nil, status.Errorf(400, "invalid product id, %s", err.Error())
	}

	o, err := db.GetClient().Order.Create().
		SetUUID(uuid.New()).
		SetProductID(pid).
		SetPrice(p.Price).
		SetUserID(acc.UUID).
		SetNote(req.Note).
		SetStatus(order.StatusUnpaid).
		Save(ctx)
	if err != nil {
		return nil, status.Errorf(500, "failed to create order, err: %s ", err)
	}

	dto := NewOrderDTO(o, p)

	return dto.ToProto(), nil
}

func (s *OrderServiceImpl) GetOrder(ctx context.Context, acc *ent.Account, req *model.GetOrderRequest) (*model.Order, error) {
	uid, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(400, "invalid id, %s", err.Error())
	}
	o, err := dao.GetOrder(ctx, uid, acc.UUID)
	if err != nil {
		return nil, err
	}
	p, err := dao.GetProduct(ctx, o.ProductID)
	if err != nil {
		return nil, status.Errorf(400, "invalid product id, %s", err.Error())
	}

	//nolint:errcheck
	go processOrder(o.UserID, o.UUID, o.Status)

	dto := NewOrderDTO(o, p)
	return dto.ToProto(), nil
}

func (s *OrderServiceImpl) ListOrders(ctx context.Context, acc *ent.Account, req *model.ListOrdersRequest) (*model.ListOrdersResponse, error) {
	limit, offset := int(req.Limit), int(req.Offset)
	if limit == 0 {
		limit = 10
	}
	list, err := dao.ListOrders(ctx, acc.UUID, order.Status(req.Status), limit, offset)
	logger.SugaredLogger.Debugw("list orders", "list", list, "err", err)
	if err != nil {
		return nil, err
	}
	ret := make([]*model.Order, 0, len(list))
	for _, o := range list {
		p, err := dao.GetProduct(ctx, o.ProductID)
		if err != nil {
			return nil, status.Errorf(400, "invalid product id, %s", err.Error())
		}

		dto := NewOrderDTO(o, p)
		ret = append(ret, dto.ToProto())
	}
	return &model.ListOrdersResponse{Orders: ret}, nil
}

func (s *OrderServiceImpl) DeleteOrder(ctx context.Context, acc *ent.Account, req *model.DeleteOrderRequest) (*emptypb.Empty, error) {
	uid, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(400, "invalid id, %s", err.Error())
	}
	err = dao.DeleteOrder(ctx, uid, acc.UUID)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *OrderServiceImpl) UpdateOrder(ctx context.Context, acc *ent.Account, req *model.UpdateOrderRequest) (*model.Order, error) {
	uid, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(400, "invalid id, %s", err.Error())
	}
	// update order
	if req.Note == "" {
		return nil, status.Errorf(400, "invalid note")
	}
	err = db.GetClient().Order.Update().Where(order.UUIDEQ(uid)).SetNote(req.Note).Exec(ctx)
	if err != nil {
		return nil, err
	}
	// retrieve order
	o, err := dao.GetOrder(ctx, uid, acc.UUID)
	if err != nil {
		return nil, err
	}
	p, err := dao.GetProduct(ctx, o.ProductID)
	if err != nil {
		return nil, err
	}
	// return order
	dto := NewOrderDTO(o, p)
	return dto.ToProto(), nil
}

// PayOrder generate qrcode for wechat pay
func (s *OrderServiceImpl) PayOrder(ctx context.Context, acc *ent.Account, req *model.PayOrderRequest) (*model.PayOrderResponse, error) {
	uid, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(400, "invalid id, %s", err.Error())
	}
	o, err := dao.GetOrder(ctx, uid, acc.UUID)
	if err != nil {
		return nil, err
	}

	if acc.UUID != o.UserID {
		return nil, status.Errorf(403, "forbidden")
	}

	p, err := dao.GetProduct(ctx, o.ProductID)
	if err != nil {
		return nil, status.Errorf(400, "invalid product id, %s", err.Error())
	}

	dto := NewOrderDTO(o, p)
	res := &model.PayOrderResponse{
		Status: "ok",
	}

	ua := ctx.Value(config.ContextKeyUserAgent)
	if ua != nil && strings.Contains(ua.(string), "MicroMessenger") {
		logger.SugaredLogger.Infow("PayOrder in wechat", "order", dto)
		// 微信浏览器内支付，使用 JSAPI
		prePayId, err := prepayInWexin(ctx, dto, acc.UserID)
		if err != nil {
			logger.SugaredLogger.Errorw("PayOrder prepay error", "err", err, "order", dto)
			return nil, err
		}
		bConfig, err := BuildJSAPIBraigeConfig(prePayId)
		if err != nil {
			logger.SugaredLogger.Errorw("PayOrder prepay error", "err", err, "order", dto)
			return nil, err
		}
		res.Config = bConfig
	} else {
		qrcode, err := Prepay(ctx, dto)
		if err != nil {
			logger.SugaredLogger.Errorw("PayOrder prepay error", "err", err, "order", dto)
			return nil, err
		}
		res.Qrcode = qrcode
	}
	err = dao.UpdateOrderStatus(db.GetClient(), ctx, o.UUID, acc.UUID, order.StatusPaying)
	if err != nil {
		logger.SugaredLogger.Errorw("PayOrder update order status error", "err", err, "order", dto)
		return nil, err
	}
	logger.SugaredLogger.Infow("PayOrder success", "res", res)
	return res, nil
}

// WechatCallBack wechat pay callback
func (s *OrderServiceImpl) WechatPayCallBack(ctx context.Context, r *http.Request) error {
	handler, err := pay.NewNotifyHandler()
	if err != nil {
		logger.SugaredLogger.Errorw("WechatCallBack new handler error", "err", err)
	}

	transaction := new(payments.Transaction)
	req, err := handler.ParseNotifyRequest(ctx, r, transaction)
	if err != nil {
		logger.SugaredLogger.Errorw("WechatCallBack parse request error", "err", err)
		return err
	}
	logger.SugaredLogger.Infow("WechatCallBack", "summary", req.Summary, "transaction", transaction)

	userId, err := uuid.Parse(*transaction.Attach)
	if err != nil {
		logger.SugaredLogger.Errorw("WechatCallBack parse user id error", "err", err)
		return err
	}

	orderId, err := uuid.Parse(*transaction.OutTradeNo)
	if err != nil {
		logger.SugaredLogger.Errorw("WechatCallBack parse order id error", "err", err)
		return err
	}
	status := order.StatusPaying
	//nolint:errcheck
	go processOrder(userId, orderId, status)
	return nil
}
