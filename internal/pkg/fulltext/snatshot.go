package fulltext

import (
	"context"

	"notionboy/internal/pkg/config"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
)

const imageQuality = 90

func SaveSnapshot(ctx context.Context, urlStr string, tag string) ([]byte, string, error) {
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()
	var buf []byte
	var title string

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
			log.Errorf("Generate pdf snapshot error: %v", err)
		}
	} else {
		// capture entire browser viewport, returning png with quality=90
		if err := chromedp.Run(ctx, chromedp.Tasks{
			chromedp.Navigate(urlStr),
			chromedp.FullScreenshot(&buf, imageQuality),
			chromedp.Title(&title),
		}); err != nil {
			log.Errorf("Generate snapshot for url: %s error: %v", urlStr, err)
			return buf, title, err
		}
	}

	log.Debugf("Success get page: %s, from url: %s", title, urlStr)
	return buf, title, nil
}
