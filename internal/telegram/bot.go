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

/* Commands for telegram bot
apikey - Get your Notion API key
bind - Bind your Telegram account to your Notion account
chat - Chat with ChatGPT
contact - Get the contact information of the author
help - Display help message
magiccode - Get a magic code for logging into the web UI
start - Start the bot
unbind - Unbind your Telegram account from your Notion account
webui - Open the NotionBoy web page
whoami - Display information about your account
zlib - Search for Zlibrary books
*/

func (b *TeleBot) registerHandlers() {
	b.Handle("/chat", handler.OnChatGPT)
	b.Handle("/start", handler.OnStart)
	b.Handle("/help", handler.OnStart)
	b.Handle("/bind", handler.OnBind)
	b.Handle("/unbind", handler.OnUnbind)
	b.Handle("/zlib", handler.OnZlib)
	b.Handle("/webui", handler.OnWebUI)
	b.Handle("/MagicCode", handler.OnMagicCode)
	b.Handle("/whoami", handler.OnWhoAmI)
	b.Handle("/sos", handler.OnSOS)
	b.Handle("/apikey", handler.OnApiKey)

	b.Handle(&tele.InlineButton{Unique: handler.INLINE_UNIQUE_ZLIB_SEARCHER}, handler.OnZlib)
	b.Handle(&tele.InlineButton{Unique: handler.INLINE_UNIQUE_ZLIB_SAVE_TO_NOTION}, handler.OnZlibSaveToNotion)

	b.Handle(tele.OnText, handler.OnText)
	b.Handle(tele.OnMedia, handler.OnMedia)
}

func Serve() {
	b := DefaultBot()
	logger.SugaredLogger.Info("Telegram bot started")
	b.Start()
}
