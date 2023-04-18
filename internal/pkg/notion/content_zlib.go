package notion

import (
	"fmt"

	"notionboy/internal/zlib"

	"github.com/jomei/notionapi"
)

type ZlibContent struct {
	Books []*zlib.Book `json:"books"`
}

func (c *ZlibContent) BuildBlocks() []notionapi.Block {
	var blocks []notionapi.Block

	for _, book := range c.Books {
		blocks = append(blocks, buildBookBlock(book))
	}

	return blocks
}

func buildBookBlock(book *zlib.Book) *notionapi.ParagraphBlock {
	richTextBlocks := make([]notionapi.RichText, 0)

	richTextBlocks = append(richTextBlocks, notionapi.RichText{
		Type: "text",
		Text: &notionapi.Text{
			Content: fmt.Sprintf("📚 %s (%d)\n", book.Title, book.Year),
		},
		Annotations: &notionapi.Annotations{
			Bold: true,
		},
	})

	if book.Author != "" {
		richTextBlocks = append(richTextBlocks, notionapi.RichText{
			Type: "text",
			Text: &notionapi.Text{
				Content: fmt.Sprintf("✍ %s\n", book.Author),
			},
		})
	}

	if book.Publisher != "" {
		richTextBlocks = append(richTextBlocks, notionapi.RichText{
			Type: "text",
			Text: &notionapi.Text{
				Content: fmt.Sprintf("🖨 %s\n", book.Publisher),
			},
		})
	}

	richTextBlocks = append(richTextBlocks, notionapi.RichText{
		Type: "text",
		Text: &notionapi.Text{
			Content: "⬇️ ",
		},
	})

	downloadLinkBlocks := make([]notionapi.RichText, 0)
	for idx, link := range book.IpfsLinks {
		downloadLinkBlocks = append(downloadLinkBlocks, notionapi.RichText{
			Type: "text",
			Text: &notionapi.Text{
				Content: fmt.Sprintf("%d ", idx+1),
				Link: &notionapi.Link{
					Url: link,
				},
			},
		})
	}
	richTextBlocks = append(richTextBlocks, downloadLinkBlocks...)
	richTextBlocks = append(richTextBlocks, notionapi.RichText{
		Type: "text",
		Text: &notionapi.Text{
			Content: fmt.Sprintf(" (%s, %s)\n", book.Extension, book.FileSizeHuman),
		},
	})

	return &notionapi.ParagraphBlock{
		BasicBlock: notionapi.BasicBlock{
			Object: notionapi.ObjectTypeBlock,
			Type:   notionapi.BlockTypeParagraph,
		},
		Paragraph: notionapi.Paragraph{
			RichText: richTextBlocks,
		},
	}
}
