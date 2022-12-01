package r2

import (
	"context"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"
	"time"

	"github.com/go-resty/resty/v2"
)

var restyClient *resty.Client

func init() {
	restyClient = resty.New()
	restyClient.SetTimeout(60 * time.Second)
}

type R2 interface {
	Upload(ctx context.Context, name, contentType string, data []byte) (string, error)
}

type R2Client struct {
	Token string
	Url   string
}

func New() R2 {
	return NewR2Client(config.GetConfig().R2.Token, config.GetConfig().R2.Url)
}

func NewR2Client(token, url string) R2 {
	return &R2Client{
		Token: token,
		Url:   url,
	}
}

func (c *R2Client) Upload(ctx context.Context, name, contentType string, data []byte) (string, error) {
	logger.SugaredLogger.Debugf("uploading to r2: %s. name: %s, contentType: %s", c.Url, name, contentType)
	url := c.Url + "/objects/" + name + "?token=" + c.Token
	logger.SugaredLogger.Debugf("url: %s", url)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	resp, err := restyClient.R().
		// SetContext(ctx).
		SetBody(data).
		SetContentLength(true).
		SetHeader("Content-Type", contentType).
		Post(url)
	logger.SugaredLogger.Debugf("upload to R2 resp: status: %s, resp: %s", resp.Status(), resp)
	if err != nil {
		logger.SugaredLogger.Errorf("upload error: %v", err)
		return "", err
	}
	return url, nil
}
