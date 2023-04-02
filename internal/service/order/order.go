package order

import (
	"context"
	"net/http"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/service/order/pay"
	"strings"
	"time"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
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
		TimeExpire:  core.Time(time.Now().Add(time.Minute * 5)),
		NotifyUrl:   core.String("https://www.weixin.qq.com/wxpay/pay.php"),
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

func QueryOrderStatus(ctx context.Context, orderDto *OrderDTO) (*payments.Transaction, error) {
	svc := pay.NewNativeAPIService()

	cfg := config.GetConfig().Wechat
	orderId := strings.ReplaceAll(orderDto.UUID.String(), "-", "")
	req := native.QueryOrderByOutTradeNoRequest{
		Mchid:      core.String(cfg.MchID),
		OutTradeNo: core.String(orderId),
	}

	resp, res, err := svc.QueryOrderByOutTradeNo(ctx, req)
	if err != nil {
		logger.SugaredLogger.Errorw("query order status QueryOrderByOutTradeNo error", "err", err, "order", orderDto)
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
