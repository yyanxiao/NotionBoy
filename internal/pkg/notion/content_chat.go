package notion

import (
	"strings"

	"github.com/jomei/notionapi"
)

type ChatContent struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	UserID   string `json:"user_id"`
}

func (c *ChatContent) BuildBlocks() []notionapi.Block {
	blocks := make([]notionapi.Block, 0)
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
						Content: c.Question,
					},
					Annotations: &notionapi.Annotations{
						Bold: true,
					},
				},
			},
		},
	})
	if c.UserID != "" {
		blocks = append(blocks, notionapi.ParagraphBlock{
			BasicBlock: notionapi.BasicBlock{
				Object: notionapi.ObjectTypeBlock,
				Type:   notionapi.BlockTypeParagraph,
			},
			Paragraph: notionapi.Paragraph{
				RichText: []notionapi.RichText{
					{
						Mention: &notionapi.Mention{
							Type: "user",
							User: &notionapi.User{
								ID: notionapi.UserID(c.UserID),
							},
						},
					},
				},
			},
		})
	}
	blocks = formatAnswer(blocks, c.Answer)
	return blocks
}

func formatAnswer(blocks []notionapi.Block, answer string) []notionapi.Block {
	content := make([]string, 0)
	codeStart := true
	language := "plain text"
	for _, txt := range strings.Split(answer, "\n") {
		if strings.HasPrefix(txt, "```") {
			if codeStart {
				blocks = append(blocks, buildTextBlock(content))
				content = []string{}
				codeStart = false
				if txt[3:] != "" {
					language = txt[3:]
				}
			} else {
				blocks = append(blocks, buildCodeBlock(content, language))
				content = []string{}
				codeStart = true
			}
		} else {
			content = append(content, txt)
		}
	}
	if len(content) > 0 {
		if codeStart {
			blocks = append(blocks, buildTextBlock(content))
		} else {
			blocks = append(blocks, buildCodeBlock(content, language))
		}
	}
	return blocks
}

func buildTextBlock(content []string) notionapi.Block {
	return notionapi.ParagraphBlock{
		BasicBlock: notionapi.BasicBlock{
			Object: notionapi.ObjectTypeBlock,
			Type:   notionapi.BlockTypeParagraph,
		},
		Paragraph: notionapi.Paragraph{
			RichText: []notionapi.RichText{
				{
					Type: "text",
					Text: &notionapi.Text{
						Content: strings.Join(content, "\n"),
					},
				},
			},
		},
	}
}

func buildCodeBlock(content []string, language string) notionapi.Block {
	return notionapi.CodeBlock{
		BasicBlock: notionapi.BasicBlock{
			Object: notionapi.ObjectTypeBlock,
			Type:   notionapi.BlockTypeCode,
		},
		Code: notionapi.Code{
			Language: language,
			RichText: []notionapi.RichText{
				{
					Type: "text",
					Text: &notionapi.Text{
						Content: strings.Join(content, "\n"),
					},
				},
			},
		},
	}
}
