package notion

import (
	"context"
	"fmt"
	"notionboy/db/ent"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/fulltext"
	"notionboy/internal/pkg/storage"
	"strings"

	"github.com/jomei/notionapi"
)

type SnapshotContent struct {
	URL      string       `json:"url"`
	Title    string       `json:"title"`
	ImageURL string       `json:"image_url"`
	PDFURL   string       `json:"pdf_url"`
	Account  *ent.Account `json:"account"`
}

func (c *SnapshotContent) ProcessSnapshot(ctx context.Context, tag string) {
	if strings.HasPrefix(strings.ToUpper(tag), config.CMD_PDF) {
		c.processPDF(ctx)
	} else {
		c.processImage(ctx)
	}
}

func (c *SnapshotContent) processImage(ctx context.Context) {
	buf, title, err := fulltext.SaveSnapshot(ctx, c.URL, config.CMD_SNAPSHOT)
	if err != nil {
		return
	}
	c.Title = title
	key := c.buildS3Key("jpg")
	imgUrl, err := storage.DefaultClient().Upload(ctx, key, buf)
	if err != nil {
		return
	}
	c.ImageURL = imgUrl
}

func (c *SnapshotContent) processPDF(ctx context.Context) {
	buf, title, err := fulltext.SaveSnapshot(ctx, c.URL, config.CMD_SNAPSHOT_PDF)
	if err != nil {
		return
	}
	c.Title = title
	key := c.buildS3Key("pdf")
	imgUrl, err := storage.DefaultClient().Upload(ctx, key, buf)
	if err != nil {
		return
	}
	c.PDFURL = imgUrl
}

func (c *SnapshotContent) buildS3Key(ext string) string {
	return fmt.Sprintf("%s/%s/%s.%s", c.Account.UserType, c.Account.UserID, c.Title, ext)
}

func (c *SnapshotContent) BuildBlocks() []notionapi.Block {
	var imageBlock notionapi.ImageBlock
	var pdfBlock notionapi.PdfBlock
	var originMediaURL string

	blocks := make([]notionapi.Block, 0, 4)
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
						Content: "\n阅读原文\n",
						Link: &notionapi.Link{
							Url: c.URL,
						},
					},
				},
			},
		},
	})

	if c.PDFURL != "" {
		originMediaURL = c.PDFURL
		pdfBlock = notionapi.PdfBlock{
			BasicBlock: notionapi.BasicBlock{
				Object: notionapi.ObjectTypeBlock,
				Type:   notionapi.BlockTypePdf,
			},
			Pdf: notionapi.Pdf{
				Type: "external",
				External: &notionapi.FileObject{
					URL: c.PDFURL,
				},
			},
		}
		blocks = append(blocks, pdfBlock)
	}
	if c.ImageURL != "" {
		originMediaURL = c.ImageURL
		imageBlock = notionapi.ImageBlock{
			BasicBlock: notionapi.BasicBlock{
				Object: notionapi.ObjectTypeBlock,
				Type:   notionapi.BlockTypeImage,
			},
			Image: notionapi.Image{
				Type: "external",
				External: &notionapi.FileObject{
					URL: c.ImageURL,
				},
			},
		}
		blocks = append(blocks, imageBlock)
	}
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
						Content: "打开原始文件",
						Link: &notionapi.Link{
							Url: originMediaURL,
						},
					},
				},
			},
		},
	})
	return blocks
}
