package notion

import (
	"context"
	"fmt"
	"notionboy/db/ent"
	"notionboy/internal/pkg/fulltext"

	"github.com/jomei/notionapi"
)

type FulltextContent struct {
	URL          string       `json:"url"`
	Title        string       `json:"title"`
	Author       string       `json:"author"`
	Summary      string       `json:"summary"`
	PublishDate  string       `json:"publish_date"`
	Account      *ent.Account `json:"account"`
	NotionPageID string       `json:"notion_page_id"`
	ErrorMessage string       `json:"error_message"`
}

func (c *FulltextContent) ProcessFulltext(ctx context.Context) {
	resp, err := fulltext.SaveReadabeArticle(ctx, c.URL, c.Account.AccessToken, c.NotionPageID)
	if err != nil {
		c.ErrorMessage = err.Error()
		return
	}
	c.Title = resp.Title
	c.Author = resp.Author
	c.Summary = resp.Summary
	c.PublishDate = resp.PublishDate
}

func (c *FulltextContent) BuildBlocks() []notionapi.Block {
	blocks := make([]notionapi.Block, 0)

	if c.ErrorMessage != "" {
		blocks = append(blocks, notionapi.ParagraphBlock{
			BasicBlock: notionapi.BasicBlock{
				Object: notionapi.ObjectTypeBlock,
				Type:   notionapi.BlockTypeParagraph,
			},
			Paragraph: notionapi.Paragraph{
				RichText: []notionapi.RichText{
					{
						Type: "text",
						Text: &notionapi.Text{
							Content: fmt.Sprintf("获取文章失败， 失败原因: %s", c.ErrorMessage),
						},
					},
				},
			},
		})
		return blocks
	}

	// blocks = append(blocks, notionapi.ParagraphBlock{
	// 	BasicBlock: notionapi.BasicBlock{
	// 		Object: notionapi.ObjectTypeBlock,
	// 		Type:   notionapi.BlockTypeParagraph,
	// 	},
	// 	Paragraph: notionapi.Paragraph{
	// 		RichText: []notionapi.RichText{
	// 			{
	// 				Type: "text",
	// 				Text: &notionapi.Text{
	// 					Content: "\n阅读原文\n",
	// 					Link: &notionapi.Link{
	// 						Url: c.URL,
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// })
	return blocks
}
