package notion

import (
	"context"
	"fmt"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"

	"github.com/jomei/notionapi"
)

func bindNotion(ctx context.Context, token string) (string, error) {
	// 获取用户绑定的 Database ID，如果有多个，只取找到的第一个
	databaseID, err := getDatabaseID(ctx, token)
	if err != nil {
		logger.SugaredLogger.Errorw("Get notion database id error", "err", err)
		return "", err
	}

	// 第一次绑定的时候自动建立 Text 和 Tags 等 DatabaseProperties，确保绑定成功
	n := &Notion{BearerToken: token, DatabaseID: databaseID}
	msg, err := UpdateDatabaseProperties(ctx, n)
	if err != nil {
		logger.SugaredLogger.Errorw("Update database error", "err", err)
	}

	logger.SugaredLogger.Debugw("Update database properties success", "msg", msg)
	content := &Content{Text: config.MSG_WELCOME}
	msg, _, err = n.CreateRecord(ctx, content)
	logger.SugaredLogger.Infow("Create database record success", "msg", msg)
	if err != nil {
		return "", err
	}
	return databaseID, nil
}

func getDatabaseID(ctx context.Context, token string) (string, error) {
	cli := notionapi.NewClient(notionapi.Token(token), func(c *notionapi.Client) {})
	searchFilter := make(map[string]string)
	searchFilter["property"] = "object"
	searchFilter["value"] = "database"
	searchReq := notionapi.SearchRequest{
		PageSize: 1,
		Filter: map[string]string{
			"property": "object",
			"value":    "database",
		},
	}
	res, err := cli.Search.Do(ctx, &searchReq)
	if err != nil {
		logger.SugaredLogger.Errorw("Search database error", "err", err)
		return "", err
	}
	databases := res.Results
	if len(databases) == 0 {
		return "", fmt.Errorf("至少需要绑定一个 Database")
	}
	database := databases[0].(*notionapi.Database)

	logger.SugaredLogger.Debugw("Found database", "database", database)
	databaseId := database.ID.String()
	return databaseId, nil
}

func UpdateDatabaseProperties(ctx context.Context, cfg *Notion) (string, error) {
	return cfg.UpdateDatabase(ctx, defaultDatabaseProperties())
}

func defaultDatabaseProperties() *notionapi.DatabaseUpdateRequest {
	return &notionapi.DatabaseUpdateRequest{
		Properties: notionapi.PropertyConfigs{
			"Tags": notionapi.MultiSelectPropertyConfig{
				Type: notionapi.PropertyConfigTypeMultiSelect,
				MultiSelect: notionapi.Select{
					Options: []notionapi.Option{},
				},
			},
			"Text": notionapi.RichTextPropertyConfig{
				Type: notionapi.PropertyConfigTypeRichText,
			},
			"Name": notionapi.TitlePropertyConfig{
				Type: notionapi.PropertyConfigTypeTitle,
			},
			"CreatedAt": notionapi.CreatedTimePropertyConfig{
				Type: notionapi.PropertyConfigCreatedTime,
			},
			"UpdatedAt": notionapi.LastEditedTimePropertyConfig{
				Type: notionapi.PropertyConfigLastEditedTime,
			},
			"CreatedBy": notionapi.CreatedByPropertyConfig{
				Type: notionapi.PropertyConfigCreatedBy,
			},
			"UpdatedBy": notionapi.LastEditedByPropertyConfig{
				Type: notionapi.PropertyConfigLastEditedBy,
			},
		},
	}
}
