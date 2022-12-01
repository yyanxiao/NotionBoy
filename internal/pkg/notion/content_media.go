package notion

import (
	"strings"

	"github.com/jomei/notionapi"
)

type MediaContent struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

func (c *MediaContent) BuildBlocks() []notionapi.Block {
	mediaType := c.Type
	mediaURL := c.URL
	var mediaBlock notionapi.Block
	if strings.HasPrefix(mediaType, "image") {
		mediaBlock = notionapi.ImageBlock{
			BasicBlock: notionapi.BasicBlock{
				Object: notionapi.ObjectTypeBlock,
				Type:   notionapi.BlockTypeImage,
			},
			Image: notionapi.Image{
				Type: "external",
				External: &notionapi.FileObject{
					URL: mediaURL,
				},
			},
		}
	} else if strings.HasPrefix(mediaType, "video") {
		mediaBlock = notionapi.VideoBlock{
			BasicBlock: notionapi.BasicBlock{
				Object: notionapi.ObjectTypeBlock,
				Type:   notionapi.BlockTypeVideo,
			},
			Video: notionapi.Video{
				Type: "external",
				External: &notionapi.FileObject{
					URL: mediaURL,
				},
			},
		}
	} else {
		mediaBlock = notionapi.FileBlock{
			BasicBlock: notionapi.BasicBlock{
				Object: notionapi.ObjectTypeBlock,
				Type:   notionapi.BlockTypeFile,
			},
			File: notionapi.BlockFile{
				Type: "external",
				External: &notionapi.FileObject{
					URL: mediaURL,
				},
			},
		}
	}

	return []notionapi.Block{
		mediaBlock,
		notionapi.ParagraphBlock{
			BasicBlock: notionapi.BasicBlock{
				Object: notionapi.ObjectTypeBlock,
				Type:   notionapi.BlockTypeParagraph,
			},
			Paragraph: notionapi.Paragraph{
				RichText: []notionapi.RichText{
					{
						Type: "text",
						Text: &notionapi.Text{
							Content: mediaURL,
							Link: &notionapi.Link{
								Url: mediaURL,
							},
						},
					},
				},
			},
		},
	}
}
