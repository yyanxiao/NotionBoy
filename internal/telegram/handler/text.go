package handler

import (
	"context"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/pkg/notion"

	tele "gopkg.in/telebot.v3"
)

func OnText(c tele.Context) error {
	ctx := context.Background()
	acc, err := queryUserAccount(ctx, c)
	if err != nil {
		logger.SugaredLogger.Errorf("Query Account Error: %v", err)
		return c.Reply(config.MSG_ERROR_ACCOUNT_NOT_FOUND)
	}

	nContent := &notion.Content{
		Text: c.Message().Text,
	}

	nt := &notion.Notion{BearerToken: acc.AccessToken, DatabaseID: acc.DatabaseID}
	res, pageID, err := nt.CreateRecord(ctx, &notion.Content{
		Text: "内容正在更新，请稍等",
	})
	nContent.Process(ctx)
	if err == nil {
		nt.PageID = pageID
		nt.UpdateRecord(ctx, nContent)
	}

	return c.Reply(res)
}
