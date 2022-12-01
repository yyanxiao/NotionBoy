package fulltext

import (
	"context"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

const imageQuality = 90

func SaveSnapshot(ctx context.Context, urlStr string, tag string) ([]byte, string, error) {
	// use remote devtools if available, else use the localhost
	if config.GetConfig().DevToolsURL != "" {
		allocatorContext, cancelRemote := chromedp.NewRemoteAllocator(ctx, config.GetConfig().DevToolsURL)
		defer cancelRemote()
		ctx = allocatorContext
	}

	var buf []byte
	var title string

	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	if tag == config.CMD_FULLTEXT_PDF {
		if err := chromedp.Run(ctx, chromedp.Tasks{
			chromedp.Navigate(urlStr),
			chromedp.ActionFunc(func(ctx context.Context) error {
				pdf, _, err := page.PrintToPDF().WithPrintBackground(false).Do(ctx)
				if err != nil {
					return err
				}
				buf = pdf
				return nil
			}),
			chromedp.Title(&title),
		}); err != nil {
			logger.SugaredLogger.Errorw("Generate pdf snapshot error", "err", err)
		}
	} else {
		// capture entire browser viewport, returning png with quality=90
		if err := chromedp.Run(ctx, chromedp.Tasks{
			chromedp.Navigate(urlStr),
			chromedp.FullScreenshot(&buf, imageQuality),
			chromedp.Title(&title),
		}); err != nil {
			logger.SugaredLogger.Errorw("Generate image snapshot error", "err", err)

			return buf, title, err
		}
	}

	logger.SugaredLogger.Debugw("Success get page", "url", urlStr, "title", title)

	return buf, title, nil
}
