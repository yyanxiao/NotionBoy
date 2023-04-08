package pay

import (
	"context"

	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

var (
	client  *core.Client
	handler *notify.Handler
)

func NewClient() *core.Client {
	if client == nil {
		var err error
		cfg := config.GetConfig().Wechat

		// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
		mchPrivateKey, err := utils.LoadPrivateKeyWithPath(cfg.MchPrivateKeyPath)
		if err != nil {
			logger.SugaredLogger.Panicw("load merchant private key error", "err", err)
		}

		ctx := context.Background()
		// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
		opts := []core.ClientOption{
			option.WithWechatPayAutoAuthCipher(cfg.MchID, cfg.MchCertificateSerialNumber, mchPrivateKey, cfg.MchAPIv3Key),
		}
		client, err = core.NewClient(ctx, opts...)
		if err != nil {
			logger.SugaredLogger.Panicw("new wechat pay client error", "err", err)
		}
	}
	return client
}

func NewNativeAPIService() native.NativeApiService {
	client := NewClient()
	return native.NativeApiService{Client: client}
}

func NewJSAPIService() jsapi.JsapiApiService {
	client := NewClient()
	return jsapi.JsapiApiService{Client: client}
}

func NewNotifyHandler() (*notify.Handler, error) {
	var err error
	if handler == nil {
		handler, err = newHandler()
		if err != nil {
			logger.SugaredLogger.Panicw("new wechat pay handler error", "err", err)
		}
	}
	return handler, nil
}

func newHandler() (*notify.Handler, error) {
	ctx := context.Background()
	cfg := config.GetConfig().Wechat
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(cfg.MchPrivateKeyPath)
	if err != nil {
		logger.SugaredLogger.Panicw("load merchant private key error", "err", err)
		return nil, err
	}
	d, err := downloader.NewCertificateDownloader(ctx, cfg.MchID, mchPrivateKey, cfg.MchCertificateSerialNumber, cfg.MchAPIv3Key)
	if err != nil {
		logger.SugaredLogger.Errorw("WechatCallBack new downloader error", "err", err)
		return nil, err
	}

	handler, err := notify.NewRSANotifyHandler(config.GetConfig().Wechat.MchAPIv3Key, verifiers.NewSHA256WithRSAVerifier(d))
	if err != nil {
		logger.SugaredLogger.Errorw("WechatCallBack new handler error", "err", err)
		return nil, err
	}
	return handler, nil
}
