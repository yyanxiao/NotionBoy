package fulltext

import (
	"context"
	"fmt"
	"net/http"

	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"

	"github.com/go-resty/resty/v2"
)

var client *resty.Client

func init() {
	client = resty.New()
}

type ParseResult struct {
	Title         string `json:"title"`
	Content       string `json:"content"`
	Author        string `json:"author"`
	PublishDate   string `json:"date_published"`
	LeadImageUrl  string `json:"lead_image_url"`
	Dek           string `json:"dek"`
	NextPageUrl   string `json:"next_page_url"`
	Url           string `json:"url"`
	Domain        string `json:"domain"`
	Summary       string `json:"excerpt"` // summary of the article
	WordCount     int    `json:"word_count"`
	Direction     string `json:"direction"`
	TotalPages    int    `json:"total_pages"`
	RenderedPages int    `json:"rendered_pages"`
}

func SaveReadabeArticle(ctx context.Context, articleUrl, notionToken, notionPageID string) (*ParseResult, error) {
	url := config.GetConfig().Readability.Host + "/api/parser_to_notion"
	var result ParseResult
	resp, err := client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"url":             articleUrl,
			"notion_token":    notionToken,
			"notion_page_id":  notionPageID,
			"is_upload_image": "false",
		}).
		SetResult(&result).
		Get(url)
	if err != nil {
		logger.SugaredLogger.Errorw("Save article failed", "err", err, "url", articleUrl)
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		logger.SugaredLogger.Errorw("Save article failed",
			"status_code", resp.StatusCode,
			"url", articleUrl,
			"response", resp.String(),
		)
		return nil, fmt.Errorf("Save article failed, status code: %d; error: %s", resp.StatusCode(), resp.String())
	}

	return &result, nil
}
