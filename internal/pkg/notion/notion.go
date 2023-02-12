package notion

import (
	"context"
	"fmt"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"
	"strings"
	"sync"

	"github.com/jomei/notionapi"
)

const maxRetryCnt = 3

type Notion struct {
	DatabaseID  string            `json:"database_id"`
	BearerToken string            `json:"bearer_token"`
	PageID      string            `json:"page_id"`
	Client      *notionapi.Client `json:"client"`
}

func (n *Notion) GetClient() *notionapi.Client {
	if n.Client == nil {
		n.Client = notionapi.NewClient(notionapi.Token(n.BearerToken), func(c *notionapi.Client) {})
	}
	return n.Client
}

func (n *Notion) CreateRecord(ctx context.Context, content *Content) (string, string, error) {
	pageProperties := content.BuildPageProperties()
	blocks := content.BuildBlocks()
	pageCreateRequest := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			DatabaseID: notionapi.DatabaseID(n.DatabaseID),
		},
		Properties: *pageProperties,
		Children:   blocks,
	}

	page, err := n.GetClient().Page.Create(ctx, pageCreateRequest)
	msg := ""
	if content.Account != nil && content.Account.NotionUserID == "" {
		msg += config.MSG_USING_NOTION_TEST_ACCOUNT
	}
	var pageID string
	if err != nil {
		msg += fmt.Sprintf("创建 Note 失败，失败原因, %v", err)
		logger.SugaredLogger.Error(msg)
	} else {
		pageID = strings.Replace(page.ID.String(), "-", "", -1)
		msg += fmt.Sprintf("创建 Note 成功，如需编辑更多，请前往 https://www.notion.so/%s", pageID)
		logger.SugaredLogger.Info(msg)
	}
	return msg, pageID, err
}

func (n *Notion) UpdateRecord(ctx context.Context, content *Content) {
	client := n.GetClient()
	var wg sync.WaitGroup
	wg.Add(2)

	update := func(wg *sync.WaitGroup, updateFunc func(ctx context.Context, client *notionapi.Client, pageID string, content *Content) (string, error)) {
		var err error
		for i := 0; i < maxRetryCnt; i++ {
			if _, err = updateFunc(ctx, client, n.PageID, content); err == nil {
				break
			}
		}
		if err != nil {
			updateErrorBlock(ctx, client, n.PageID, err)
		}
		wg.Done()
	}

	go update(&wg, updatePage)
	go update(&wg, updateBlock)

	wg.Wait()
}

func updateErrorBlock(ctx context.Context, client *notionapi.Client, pageID string, err error) {
	errorMsgBlock := notionapi.ParagraphBlock{
		BasicBlock: notionapi.BasicBlock{
			Object: notionapi.ObjectTypeBlock,
			Type:   notionapi.BlockTypeParagraph,
		},
		Paragraph: notionapi.Paragraph{
			RichText: []notionapi.RichText{
				{
					Type: "text",
					Text: &notionapi.Text{
						Content: "Save Content error, error message: " + err.Error(),
					},
					Annotations: &notionapi.Annotations{
						Bold: true,
					},
				},
			},
		},
	}
	_, _ = client.Block.AppendChildren(ctx, notionapi.BlockID(pageID), &notionapi.AppendBlockChildrenRequest{
		Children: []notionapi.Block{errorMsgBlock},
	})
}

func updatePage(ctx context.Context, client *notionapi.Client, pageID string, content *Content) (string, error) {
	pageProperties := content.BuildPageProperties()
	pageUpdateRequest := &notionapi.PageUpdateRequest{
		Properties: *pageProperties,
	}

	_, err := client.Page.Update(ctx, notionapi.PageID(pageID), pageUpdateRequest)

	msg := ""
	if content.Account.NotionUserID == "" {
		msg += config.MSG_USING_NOTION_TEST_ACCOUNT
	}
	if err != nil {
		msg += fmt.Sprintf("更新 Page(%s) 失败，失败原因, %v", pageID, err)
		logger.SugaredLogger.Error(msg)
	} else {
		msg += fmt.Sprintf("更新 Page(%s) 成功，如需编辑更多，请前往 https://www.notion.so/%s", pageID, pageID)
		logger.SugaredLogger.Info(msg)
	}
	return msg, err
}

func updateBlock(ctx context.Context, client *notionapi.Client, pageID string, content *Content) (string, error) {
	blocks := content.BuildBlocks()
	appendBlockChildrenRequest := &notionapi.AppendBlockChildrenRequest{
		Children: blocks,
	}
	_, err := client.Block.AppendChildren(ctx, notionapi.BlockID(pageID), appendBlockChildrenRequest)
	msg := ""
	if content.Account.NotionUserID == "" {
		msg += config.MSG_USING_NOTION_TEST_ACCOUNT
	}
	if err != nil {
		msg += fmt.Sprintf("更新 Blocks in Page(%s) 失败，失败原因, %v", pageID, err)
		logger.SugaredLogger.Error(msg)
	} else {
		msg += fmt.Sprintf("更新 Blocks in Page(%s) 成功，如需编辑更多，请前往 https://www.notion.so/%s", pageID, pageID)
		logger.SugaredLogger.Info(msg)
	}
	return msg, err
}

func (n *Notion) UpdateDatabase(ctx context.Context, req *notionapi.DatabaseUpdateRequest) (string, error) {
	_, err := n.GetClient().Database.Update(ctx, notionapi.DatabaseID(n.DatabaseID), req)
	var msg string
	if err != nil {
		msg = fmt.Sprintf("Update Database(%s) 失败，失败原因, %v", n.DatabaseID, err)
		logger.SugaredLogger.Error(msg)
	} else {
		msg = fmt.Sprintf("成功更新 Database(%s)", n.DatabaseID)
	}
	return msg, err
}

func (n *Notion) GetDatabaseInfo(ctx context.Context) (*notionapi.Database, error) {
	return n.GetClient().Database.Get(ctx, notionapi.DatabaseID(n.DatabaseID))
}
