package wechat

import (
	"context"
	"notionboy/internal/pkg/logger"

	"github.com/eatmoreapple/openwechat"
)

type WechatBot struct {
	*openwechat.Bot
	ctx context.Context
}

func New() *WechatBot {
	bot := &WechatBot{
		Bot: openwechat.DefaultBot(openwechat.Desktop),
		ctx: context.Background(),
	}
	bot.registerHandlers()
	bot.registerCallbacks()
	return bot
}

func (b *WechatBot) Serve() {
	defer func() {
		if v := recover(); v != nil {
			logger.SugaredLogger.Errorf("Recover from bot panic, err: %#v", v)
		}
	}()
	if err := b.Block(); err != nil {
		logger.SugaredLogger.Error("Start wechat bot error")
	}
}

func (b *WechatBot) registerHandlers() {
	dispatcher := openwechat.NewMessageMatchDispatcher()
	dispatcher.OnText(chat)
	b.MessageHandler = openwechat.DispatchMessage(dispatcher)
}

func (b *WechatBot) registerCallbacks() {
	b.UUIDCallback = openwechat.PrintlnQrcodeUrl
	b.LoginCallBack = loginCallBack
}

// Serve run wechat bot, login and block
func Serve() {
	b := New()
	if err := b.login(""); err != nil {
		logger.SugaredLogger.Errorw("Login to wechat error", "err", err)
		return
	}
	b.Serve()
}

// ReServe after main groutine exit
// we can use ReServe to manually login and block
func ReServe() {
	b := New()
	uuid, err := b.Caller.GetLoginUUID()
	if err != nil {
		logger.SugaredLogger.Errorw("Generate uuid for wechat error", "err", err)
	}
	if err := b.login(uuid); err != nil {
		logger.SugaredLogger.Errorw("Login to wechat error", "err", err)
		return
	}
	b.Serve()
}
