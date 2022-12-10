package fulltext

import (
	"context"
	"errors"
	"io"
	url2 "net/url"
	"notionboy/internal/pkg/browser"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

const (
	wxURL       = "mp.weixin.qq.com"
	pageTimeout = 60
)

func SaveSnapshot(ctx context.Context, urlStr string, tag string) ([]byte, string, error) {
	var buf []byte
	var title string
	url, err := url2.Parse(urlStr)
	if err != nil {
		logger.SugaredLogger.Errorw("Invalid url", "err", err, "url", urlStr)
		return buf, title, err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		time.Sleep(pageTimeout * time.Second)
		cancel()
	}()

	p := browser.New().MustConnect().MustPage(urlStr).Context(ctx)
	p.MustWaitLoad()
	// wait body to show
	p.MustElement("body").MustWaitLoad()
	// make sure title show up
	title = p.MustEval("() => document.title").Str()
	if title == "" {
		time.Sleep(time.Second)
		title = p.MustEval("() => document.title").Str()
		if title == "" {
			return buf, title, errors.New("can not load page")
		}
	}

	// get height and scroll to bottom
	// since wechat need scroll to show images
	if url.Host == wxURL {
		pageHeight := p.MustEval("() => document.body.scrollHeight")
		scrollStepsNum := 1000
		_ = p.Mouse.Scroll(0, pageHeight.Num(), scrollStepsNum)
	}

	if tag == config.CMD_FULLTEXT_PDF {
		buf, err = generatePDF(p)
	} else {
		buf, err = generateScreenshot(p)
	}
	if err != nil {
		return buf, title, err
	}

	logger.SugaredLogger.Debugw("Success get page", "url", urlStr, "title", title)
	return buf, title, nil
}

func generatePDF(p *rod.Page) ([]byte, error) {
	var err error
	var r *rod.StreamReader
	var buf []byte
	r, err = p.PDF(&proto.PagePrintToPDF{
		DisplayHeaderFooter: true,
		HeaderTemplate:      "<div></div>",
		FooterTemplate:      `<p style="text-align:left;font-size:500%;color:blue;">Powered by NotionBoy (微信搜索 NotionBoy 关注)</p>`,
	})
	if err != nil {
		logger.SugaredLogger.Errorw("Generate pdf snapshot error", "err", err)
		return buf, err
	}
	buf, err = io.ReadAll(r)
	if err != nil {
		logger.SugaredLogger.Errorw("Generate pdf snapshot error", "err", err)
		return buf, err
	}
	return buf, err
}

func generateScreenshot(p *rod.Page) ([]byte, error) {
	buf, err := p.Screenshot(true, &proto.PageCaptureScreenshot{})
	if err != nil {
		logger.SugaredLogger.Errorw("Generate image snapshot error", "err", err)
	}
	return buf, err
}
