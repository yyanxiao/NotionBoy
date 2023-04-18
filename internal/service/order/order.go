package order

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"notionboy/api/pb/model"
	"notionboy/db/ent"
	"notionboy/db/ent/order"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/db"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/service/order/pay"

	"github.com/google/uuid"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

// Prepay 微信支付预下单，会返回一个可用于支付的二维码链接
func Prepay(ctx context.Context, orderDto *OrderDTO) (string, error) {
	cfg := config.GetConfig().Wechat
	// expireTime := time.Now().Add(time.Minute * 5).Format(time.RFC3339)
	orderId := strings.ReplaceAll(orderDto.UUID.String(), "-", "")
	req := native.PrepayRequest{
		Appid:       core.String(cfg.AppID),
		Mchid:       core.String(cfg.MchID),
		Description: core.String(orderDto.Product.Name),
		OutTradeNo:  core.String(orderId),
		Attach:      core.String(string(orderDto.UserID.String())),
		TimeExpire:  core.Time(time.Now().Add(time.Minute * 5)),
		NotifyUrl:   core.String(config.GetConfig().Service.URL + "/v1/wechatpay/callback"),
		Amount: &native.Amount{
			Currency: core.String("CNY"),
			// 订单总金额，单位为分。
			Total: core.Int64(int64(orderDto.Price * 100)),
		},
	}

	svc := pay.NewNativeAPIService()

	resp, res, err := svc.Prepay(ctx, req)
	if err != nil {
		logger.SugaredLogger.Errorw("prepay error", "err", err)
		return "", err
	} else {
		if res.Response.StatusCode != http.StatusOK {
			var msg string
			if _, err := res.Response.Body.Read([]byte(msg)); err != nil {
				logger.SugaredLogger.Errorw("prepay error, read response message error", "err", err)
			}
			logger.SugaredLogger.Errorw("prepay error", "err", err, "msg", msg)
			return "", err
		}
		return *resp.CodeUrl, nil
	}
}

func prepayInWexin(ctx context.Context, orderDto *OrderDTO, userId string) (string, error) {
	cfg := config.GetConfig().Wechat
	orderId := strings.ReplaceAll(orderDto.UUID.String(), "-", "")
	req := jsapi.PrepayRequest{
		Appid:       core.String(cfg.AppID),
		Mchid:       core.String(cfg.MchID),
		Description: core.String(orderDto.Product.Name),
		OutTradeNo:  core.String(orderId),
		TimeExpire:  core.Time(time.Now().Add(time.Minute * 5)),
		Attach:      core.String(string(orderDto.UserID.String())),
		NotifyUrl:   core.String(config.GetConfig().Service.URL + "/v1/wechatpay/callback"),
		Amount: &jsapi.Amount{
			Currency: core.String("CNY"),
			// 订单总金额，单位为分。
			Total: core.Int64(int64(orderDto.Price * 100)),
		},
		Payer: &jsapi.Payer{
			Openid: core.String(userId),
		},
	}
	svc := pay.NewJSAPIService()
	resp, res, err := svc.Prepay(ctx, req)
	if err != nil {
		logger.SugaredLogger.Errorw("prepay error", "err", err)
		return "", err
	} else {
		if res.Response.StatusCode != http.StatusOK {
			var msg string
			if _, err := res.Response.Body.Read([]byte(msg)); err != nil {
				logger.SugaredLogger.Errorw("prepay error, read response message error", "err", err)
			}
			logger.SugaredLogger.Errorw("prepay error", "err", err, "msg", msg)
			return "", err
		}
		return *resp.PrepayId, nil
	}
}

func BuildJSAPIBraigeConfig(prepayId string) (*model.PayOrderConfig, error) {
	cfg := config.GetConfig().Wechat
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(cfg.MchPrivateKeyPath)
	if err != nil {
		logger.SugaredLogger.Panicw("load merchant private key error", "err", err)
		return nil, err
	}

	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	appId := cfg.AppID
	nonceStr := strings.ReplaceAll(uuid.New().String(), "-", "")
	source := fmt.Sprintf("%s\n%s\n%s\nprepay_id=%s\n", appId, timeStamp, nonceStr, prepayId)
	paySign, err := utils.SignSHA256WithRSA(source, mchPrivateKey)
	if err != nil {
		logger.SugaredLogger.Panicw("sign error", "err", err)
		return nil, err
	}

	return &model.PayOrderConfig{
		AppId:     appId,
		Timestamp: timeStamp,
		NonceStr:  nonceStr,
		Package:   "prepay_id=" + prepayId,
		SignType:  "RSA",
		PaySign:   paySign,
	}, nil
}

// processOrder after receive wechat pay notify, will call this method to process order
func processOrder(userId, orderId uuid.UUID, status order.Status) error {
	ctx := context.Background()
	if status == order.StatusPaid {
		return processAfterPaid(ctx, userId, orderId)
	}
	return processAfterPaying(ctx, userId, orderId)
}

func processAfterPaying(ctx context.Context, userId, orderId uuid.UUID) error {
	tx, err := db.GetClient().Tx(ctx)
	if err != nil {
		logger.SugaredLogger.Errorw("processAfterPaying get tx error", "err", err, "orderId", orderId, "userId", userId)
		return err
	}

	// check order status
	status, err := queryOrderStatusForUpdate(ctx, tx.Client(), userId, orderId)
	if err != nil {
		logger.SugaredLogger.Errorw("processAfterPaying query order status error", "err", err, "userId", userId, "orderId", orderId)
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	// make sure order status is paid, then process
	if *status != order.StatusPaying {
		logger.SugaredLogger.Infow("processAfterPaying order status is not paying", "userId", userId, "orderId", orderId, "status", status)
		return tx.Rollback()
	}

	// query order status from wechat
	res, err := queryOrderPaymentStatusFromWechat(ctx, orderId.String())
	if err != nil {
		logger.SugaredLogger.Errorw("processAfterPaying query order status from wechat error", "err", err, "userId", userId, "orderId", orderId)
		return err
	}
	if *res.TradeState == "SUCCESS" {
		// update order success status and payment info
		paymentInfo, err := json.Marshal(res)
		if err != nil {
			logger.SugaredLogger.Errorw("processAfterPaying marshal payment info error", "err", err, "userId", userId, "orderId", orderId)
			return err
		}
		err = tx.Order.Update().Where(order.UUIDEQ(orderId), order.UserIDEQ(userId)).SetStatus(order.StatusPaid).SetPaymentInfo(string(paymentInfo)).Exec(ctx)
		if err != nil {
			logger.SugaredLogger.Errorw("processAfterPaying update order status to StatusPaid error", "err", err, "userId", userId, "orderId", orderId)
			return err
		}
	}
	return tx.Commit()
}

func processAfterPaid(ctx context.Context, userId, orderId uuid.UUID) error {
	tx, err := db.GetClient().Tx(ctx)
	if err != nil {
		logger.SugaredLogger.Errorw("processAfterPaid get tx error", "err", err, "userId", userId, "orderId", orderId)
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	// check order status
	status, err := queryOrderStatusForUpdate(ctx, tx.Client(), userId, orderId)
	if err != nil {
		logger.SugaredLogger.Errorw("processAfterPaid query order status error", "err", err, "userId", userId, "orderId", orderId)
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	// make sure order status is paid, then process
	if *status != order.StatusPaid {
		logger.SugaredLogger.Errorw("processAfterPaid order status is not paid", "userId", userId, "orderId", orderId, "status", status)
		return tx.Rollback()
	}

	// retrieve order
	o, err := dao.GetOrder(ctx, orderId, userId)
	if err != nil {
		logger.SugaredLogger.Errorw("processAfterPaid get order error", "err", err, "userId", userId, "orderId", orderId)
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	p, err := dao.GetProduct(ctx, o.ProductID)
	if err != nil {
		logger.SugaredLogger.Errorw("processAfterPaid get product error", "err", err, "userId", userId, "orderId", orderId)
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	acc, err := dao.QueryAccountByUUID(ctx, userId)
	if err != nil {
		logger.SugaredLogger.Errorw("processAfterPaid get account error", "err", err, "userId", userId, "orderId", orderId)
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	// return order
	dto := NewOrderDTO(o, p)

	err = dao.UpdateOrderStatus(tx.Client(), ctx, dto.UUID, acc.UUID, order.StatusProcessing)
	if err != nil {
		logger.SugaredLogger.Errorw("processAfterPaid update order status to StatusProcessing error", "err", err, "order", dto)
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	// update quota token
	err = dao.UpdateQuota(tx.Client(), ctx, acc.ID, dto.Product.Token, dto.Product.Name)
	if err != nil {
		logger.SugaredLogger.Errorw("processAfterPaid update quota error", "err", err, "order", dto)
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	// update order status to completed
	err = dao.UpdateOrderStatus(tx.Client(), ctx, dto.UUID, acc.UUID, order.StatusCompleted)
	if err != nil {
		logger.SugaredLogger.Errorw("processAfterPaid update order status error", "err", err, "order", dto, "status", order.StatusCompleted)
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	logger.SugaredLogger.Infow("processAfterPaid finished", "order", dto)
	err = tx.Commit()
	if err != nil {
		logger.SugaredLogger.Errorw("processAfterPaid commit tx error", "err", err, "order", dto)
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	return nil
}

func queryOrderStatusForUpdate(ctx context.Context, cli *ent.Client, userId uuid.UUID, orderId uuid.UUID) (*order.Status, error) {
	// check order status
	rows, err := cli.QueryContext(ctx, "select status from orders where user_id = ? and uuid = ? for update", userId.String(), orderId.String())
	if err != nil {
		logger.SugaredLogger.Errorw("Query order status for update error", "err", err, "userId", userId, "orderId", orderId)
		return nil, err
	}
	defer rows.Close()
	var status order.Status
	if rows.Next() {
		err = rows.Scan(&status)
		if err != nil {
			logger.SugaredLogger.Errorw("Query order status for update scan order status error", "err", err, "userId", userId, "orderId", orderId)
			return nil, err
		}
	}
	return &status, nil
}

func queryOrderPaymentStatusFromWechat(ctx context.Context, orderId string) (*payments.Transaction, error) {
	svc := pay.NewNativeAPIService()

	cfg := config.GetConfig().Wechat
	orderId = strings.ReplaceAll(orderId, "-", "")
	req := native.QueryOrderByOutTradeNoRequest{
		Mchid:      core.String(cfg.MchID),
		OutTradeNo: core.String(orderId),
	}

	resp, res, err := svc.QueryOrderByOutTradeNo(ctx, req)
	if err != nil {
		logger.SugaredLogger.Errorw("query order status QueryOrderByOutTradeNo error", "err", err, "orderId", orderId)
		return nil, err
	} else {
		if res.Response.StatusCode != http.StatusOK {
			var msg string
			if _, err := res.Response.Body.Read([]byte(msg)); err != nil {
				logger.SugaredLogger.Errorw("query order status QueryOrderByOutTradeNo error, read response message error", "err", err)
			}
			logger.SugaredLogger.Errorw("query order status QueryOrderByOutTradeNo error", "err", err, "msg", msg)
			return nil, err
		}
		return resp, nil
	}
}
