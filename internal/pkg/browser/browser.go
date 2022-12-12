package browser

import (
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"
	"sync"

	"github.com/go-rod/rod"
)

var (
	browser *rod.Browser
	once    sync.Once
)

func New() *rod.Browser {
	once.Do(func() {
		if config.GetConfig().DevToolsURL != "" {
			url := detectURL(config.GetConfig().DevToolsURL)
			logger.SugaredLogger.Infow("Use existing browser", "dev_url", url)
			browser = rod.New().ControlURL(url).MustConnect()
		} else {
			logger.SugaredLogger.Info("Use built-in browser")
			browser = rod.New().MustConnect()
		}
		logger.SugaredLogger.Debug("Success connect browser")
	})

	return browser
}
