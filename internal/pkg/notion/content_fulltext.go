package notion

import (
	"context"
	"net/url"

	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/fulltext"
	"notionboy/internal/pkg/r2"

	"github.com/jomei/notionapi"
)

type FulltextContent struct {
	URL      string `json:"url"`
	Title    string `json:"title"`
	ImageURL string `json:"image_url"`
	PDFURL   string `json:"pdf_url"`
}

func (c *FulltextContent) ProcessSnapshot(ctx context.Context, tag string) {
	if tag == config.CMD_FULLTEXT_PDF {
		c.processFulltextPDF(ctx)
	} else {
		c.processFulltextImage(ctx)
	}
}

func (c *FulltextContent) processFulltextImage(ctx context.Context) {
	buf, title, err := fulltext.SaveSnapshot(ctx, c.URL, config.CMD_FULLTEXT)
	if err != nil {
		return
	}
	r2Client := r2.New()
	imgUrl, err := r2Client.Upload(ctx, url.QueryEscape(title)+".jpg", "image/jpeg", buf)
	if err != nil {
		return
	}
	c.ImageURL = imgUrl
	c.Title = title
}

func (c *FulltextContent) processFulltextPDF(ctx context.Context) {
	buf, title, err := fulltext.SaveSnapshot(ctx, c.URL, config.CMD_FULLTEXT_PDF)
	if err != nil {
		return
	}
	r2Client := r2.New()
	imgUrl, err := r2Client.Upload(ctx, url.QueryEscape(title)+".pdf", "application/pdf", buf)
	if err != nil {
		return
	}
	c.PDFURL = imgUrl
	c.Title = title
}

func (c *FulltextContent) BuildBlocks() []notionapi.Block {
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
	return blocks
}
