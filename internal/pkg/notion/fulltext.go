package notion

import (
	"context"
	"net/url"

	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/fulltext"
	"notionboy/internal/pkg/r2"

	"github.com/jomei/notionapi"
)

func buildFulltextContent(pageCreateRequest *notionapi.PageCreateRequest, content *Content) {
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
					Text: notionapi.Text{
						Content: "\n阅读原文\n",
						Link: &notionapi.Link{
							Url: content.Fulltext.URL,
						},
					},
				},
			},
		},
	})

	if content.Fulltext.PDFURL != "" {
		originMediaURL = content.Fulltext.PDFURL
		pdfBlock = notionapi.PdfBlock{
			BasicBlock: notionapi.BasicBlock{
				Object: notionapi.ObjectTypeBlock,
				Type:   notionapi.BlockTypePdf,
			},
			Pdf: notionapi.Pdf{
				Type: "external",
				External: &notionapi.FileObject{
					URL: content.Fulltext.PDFURL,
				},
			},
		}
		blocks = append(blocks, pdfBlock)
	}
	if content.Fulltext.ImageURL != "" {
		originMediaURL = content.Fulltext.ImageURL
		imageBlock = notionapi.ImageBlock{
			BasicBlock: notionapi.BasicBlock{
				Object: notionapi.ObjectTypeBlock,
				Type:   notionapi.BlockTypeImage,
			},
			Image: notionapi.Image{
				Type: "external",
				External: &notionapi.FileObject{
					URL: content.Fulltext.ImageURL,
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
					Text: notionapi.Text{
						Content: "打开原始文件",
						Link: &notionapi.Link{
							Url: originMediaURL,
						},
					},
				},
			},
		},
	})
	pageCreateRequest.Children = blocks
}

func (c *Content) processFulltextSnapshot(ctx context.Context, tag string) {
	if tag == config.CMD_FULLTEXT_PDF {
		c.processFulltextPDF(ctx)
	} else {
		c.processFulltextImage(ctx)
	}
}

func (c *Content) processFulltextImage(ctx context.Context) {
	buf, title, err := fulltext.SaveSnapshot(ctx, c.Fulltext.URL, config.CMD_FULLTEXT)
	if err != nil {
		return
	}
	r2Client := r2.New()
	imgUrl, err := r2Client.Upload(ctx, url.QueryEscape(title)+".jpg", "image/jpeg", buf)
	if err != nil {
		return
	}
	c.Fulltext.ImageURL = imgUrl
	c.Fulltext.Title = title
}

func (c *Content) processFulltextPDF(ctx context.Context) {
	buf, title, err := fulltext.SaveSnapshot(ctx, c.Fulltext.URL, config.CMD_FULLTEXT_PDF)
	if err != nil {
		return
	}
	r2Client := r2.New()
	imgUrl, err := r2Client.Upload(ctx, url.QueryEscape(title)+".pdf", "application/pdf", buf)
	if err != nil {
		return
	}
	c.Fulltext.PDFURL = imgUrl
	c.Fulltext.Title = title
}
