package telegram

import (
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/telegram/handler"

	tele "gopkg.in/telebot.v3"
)

type TeleBot struct {
	*tele.Bot
}

var bot *TeleBot

func New(token string) *TeleBot {
	b, err := tele.NewBot(tele.Settings{
		Token: token,
	})
	if err != nil {
		logger.SugaredLogger.Fatalf("Init telebot error", "err", err)
	}

	teleBot := &TeleBot{Bot: b}
	teleBot.registerHandlers()

	return teleBot
}

func DefaultBot() *TeleBot {
	if bot == nil {
		bot = New(config.GetConfig().Telegram.Token)
	}
	return bot
}

func (b *TeleBot) registerHandlers() {
	b.Handle("/chat", handler.OnChatGPT)
	b.Handle("/start", handler.OnStart)
	b.Handle("/help", handler.OnStart)
	b.Handle("/bind", handler.OnBind)
	b.Handle("/unbind", handler.OnUnbind)

	b.Handle(tele.OnText, handler.OnText)
	b.Handle(tele.OnMedia, handler.OnMedia)
}

func Serve() {
	b := DefaultBot()
	logger.SugaredLogger.Info("Telegram bot started")
	b.Start()
}
