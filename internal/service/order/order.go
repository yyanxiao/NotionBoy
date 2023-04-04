package order

import (
	"context"
	"fmt"
	"net/http"
	"notionboy/api/pb/model"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/service/order/pay"
	"strconv"
	"strings"
	"time"

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

func prepayInWexin(ctx context.Context, orderDto *OrderDTO, userId string) (string, error) {
	cfg := config.GetConfig().Wechat
	orderId := strings.ReplaceAll(orderDto.UUID.String(), "-", "")
	req := jsapi.PrepayRequest{
		Appid:       core.String(cfg.AppID),
		Mchid:       core.String(cfg.MchID),
		Description: core.String(orderDto.Product.Name),
		OutTradeNo:  core.String(orderId),
		TimeExpire:  core.Time(time.Now().Add(time.Minute * 5)),
		NotifyUrl:   core.String("https://www.weixin.qq.com/wxpay/pay.php"),
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
