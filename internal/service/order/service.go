package order

import (
	"context"
	"encoding/json"
	"strings"

	"notionboy/api/pb/model"
	"notionboy/db/ent"
	"notionboy/db/ent/order"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/db"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/logger"

	"github.com/google/uuid"
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
	go processOrder(acc, o, p)

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

// processOrder after receive wechat pay notify, will call this method to process order
func processOrder(acc *ent.Account, o *ent.Order, p *ent.Product) error {
	ctx := context.Background()
	dto := NewOrderDTO(o, p)
	if o.Status == order.StatusPaying {
		return processAfterPaying(ctx, acc, dto)
	} else if o.Status == order.StatusPaid {
		return processAfterPaid(ctx, acc, dto)
	}
	return nil
}

func processAfterPaying(ctx context.Context, acc *ent.Account, dto *OrderDTO) error {
	tx, err := db.GetClient().Tx(ctx)
	if err != nil {
		logger.SugaredLogger.Errorw("processAfterPaying get tx error", "err", err, "order", dto)
		return err
	}
	// query order status from wechat
	res, err := QueryOrderStatus(ctx, dto)
	if err != nil {
		logger.SugaredLogger.Errorw("processAfterPaying query order status from wechat error", "err", err, "order", dto)
		return err
	}
	if *res.TradeState == "SUCCESS" {
		// update order success status and payment info
		paymentInfo, err := json.Marshal(res)
		if err != nil {
			logger.SugaredLogger.Errorw("processAfterPaying marshal payment info error", "err", err, "order", dto)
			return err
		}
		err = tx.Order.Update().Where(order.UUIDEQ(dto.UUID), order.UserIDEQ(dto.UserID)).SetStatus(order.StatusPaid).SetPaymentInfo(string(paymentInfo)).Exec(ctx)
		if err != nil {
			logger.SugaredLogger.Errorw("processAfterPaying update order status to StatusPaid error", "err", err, "order", dto)
			return err
		}
	}
	return tx.Commit()
}

func processAfterPaid(ctx context.Context, acc *ent.Account, dto *OrderDTO) error {
	// 1. update order status to Prcesssing
	// 2. update user quota
	// 3. update order status to Completed
	tx, err := db.GetClient().Tx(ctx)
	if err != nil {
		logger.SugaredLogger.Errorw("processAfterPaid get tx error", "err", err, "order", dto)
		return err
	}

	err = dao.UpdateOrderStatus(tx.Client(), ctx, dto.UUID, acc.UUID, order.StatusProcessing)
	if err != nil {
		logger.SugaredLogger.Errorw("processAfterPaid update order status to StatusProcessing error", "err", err, "order", dto)
		return err
	}

	// update quota token
	err = dao.UpdateQuota(tx.Client(), ctx, acc.ID, dto.Product.Token, dto.Product.Name)

	if err != nil {
		logger.SugaredLogger.Errorw("processAfterPaid update quota error", "err", err, "order", dto)
		return err
	}
	// update order status to completed
	err = dao.UpdateOrderStatus(tx.Client(), ctx, dto.UUID, acc.UUID, order.StatusCompleted)
	if err != nil {
		logger.SugaredLogger.Errorw("processAfterPaid update order status error", "err", err, "order", dto, "status", order.StatusCompleted)
		return err
	}
	return tx.Commit()
}
