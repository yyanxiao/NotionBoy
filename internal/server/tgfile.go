package server

import (
	"fmt"
	"net/http"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/telegram"
	"strings"
)

func proxyTelegramFile(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/files/tg/")
	if path == "" {
		renderError(w, http.StatusBadRequest, "Bad request", nil)
		return
	}
	id := strings.Split(path, ".")[0]

	bot := telegram.DefaultBot()
	logger.SugaredLogger.Debugw("Get file from telegram", "file_id", id)
	f, err := bot.FileByID(id)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Get file by id error", err)
		return
	}
	url := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", config.GetConfig().Telegram.Token, f.FilePath)

	logger.SugaredLogger.Debugw("TG file", "file", f, "url", url)
	http.Redirect(w, r, url, http.StatusFound)
}
