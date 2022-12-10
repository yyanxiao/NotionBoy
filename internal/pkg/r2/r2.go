package r2

import (
	"context"
	"fmt"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	restyClient *resty.Client
	client      *R2Client
)

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

func DefaultClient() R2 {
	return New(config.GetConfig().R2.Token, config.GetConfig().R2.Url)
}

func New(token, url string) R2 {
	var once sync.Once
	once.Do(func() {
		client = &R2Client{
			Token: token,
			Url:   url,
		}
	})
	return client
}

func (c *R2Client) Upload(ctx context.Context, name, contentType string, data []byte) (string, error) {
	databaseID := ctx.Value(config.DATABASE_ID)
	if databaseID.(string) != "" {
		name = fmt.Sprintf("%s-%s", databaseID, name)
	}
	logger.SugaredLogger.Debugf("uploading to r2: %s. name: %s, contentType: %s", c.Url, name, contentType)
	url := c.Url + "/objects/" + name + "?token=" + c.Token
	logger.SugaredLogger.Debugf("url: %s", url)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	resp, err := restyClient.R().
		SetContext(ctx).
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
