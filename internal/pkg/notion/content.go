package notion

import (
	"context"
	"notionboy/db/ent"
	"notionboy/internal/pkg/config"
	"regexp"
	"strings"
	"time"

	"github.com/jomei/notionapi"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Content struct {
	Tags          []string        `json:"tags"`
	Text          string          `json:"text"`
	IsFulltext    bool            `json:"is_fulltext"`
	Fulltext      FulltextContent `json:"fulltext"`
	IsMedia       bool            `json:"is_media"`
	Media         MediaContent    `json:"media"`
	IsChatContent bool            `json:"is_chat_content"`
	ChatContent   ChatContent     `json:"chat_content"`
	Medias        []*MediaContent `json:"medias"`
	Account       *ent.Account    `json:"account"`
}

// Process 从 text 提取 tags，配置全文 snapshot
func (c *Content) Process(ctx context.Context) {
	r, _ := regexp.Compile(`#(.+?)($|\s)`)
	match := r.FindAllStringSubmatch(c.Text, -1)
	if len(match) > 0 {
		tags := make([]string, 0)
		for _, m := range match {
			tag := strings.Trim(m[1], "# ")
			tags = append(tags, tag)
			if strings.ToUpper(tag) == config.CMD_FULLTEXT || strings.HasPrefix(strings.ToUpper(tag), config.CMD_PDF) {
				c.parseFulltextURL(ctx, tag)
			}
		}
		c.Tags = tags
	}
}

func (c *Content) parseFulltextURL(ctx context.Context, tag string) {
	r, _ := regexp.Compile(`https?://(www.)?[-a-zA-Z0-9@:%._+~#=]{2,256}.[a-z]{2,4}\b([-a-zA-Z0-9@:%_+.~#?&//=]*)`)
	match := r.FindAllStringSubmatch(c.Text, -1)
	if len(match) > 0 {
		// only save last url
		for _, m := range match {
			c.Fulltext.URL = m[0]
			c.Fulltext.Account = c.Account
			c.IsFulltext = true
		}
	}

	if c.IsFulltext {
		c.Fulltext.ProcessSnapshot(ctx, tag)
	}
}

func (c *Content) buildTitle() string {
	title := c.Text
	if c.IsFulltext && c.Fulltext.Title != "" {
		title = c.Fulltext.Title
	}
	if c.IsMedia && c.Text == "" {
		loc, _ := time.LoadLocation("Asia/Shanghai")
		title = cases.Title(language.English).String(c.Media.Type) + " " + time.Now().UTC().In(loc).Format(time.RFC3339)
	}
	return title
}

func (c *Content) buildTagsProperties() notionapi.MultiSelectProperty {
	var multiSelect []notionapi.Option
	for _, tag := range c.Tags {
		selectOption := notionapi.Option{
			Name: tag,
		}
		multiSelect = append(multiSelect, selectOption)
	}
	return notionapi.MultiSelectProperty{
		Type:        "multi_select",
		MultiSelect: multiSelect,
	}
}

func (c *Content) BuildBlocks() []notionapi.Block {
	blocks := make([]notionapi.Block, 0)
	if c.IsFulltext {
		fulltextBlocks := c.Fulltext.BuildBlocks()
		blocks = append(blocks, fulltextBlocks...)
	}
	if c.IsMedia {
		mediaBlocks := c.Media.BuildBlocks()
		blocks = append(blocks, mediaBlocks...)
	}
	if c.IsChatContent {
		blocks = c.ChatContent.BuildBlocks()
	}
	for _, media := range c.Medias {
		mBlocks := media.BuildBlocks()
		blocks = append(blocks, mBlocks...)
	}
	return blocks
}

func (c *Content) BuildPageProperties() *notionapi.Properties {
	pageProperties := notionapi.Properties{
		"Name": notionapi.TitleProperty{
			Type: "title",
			Title: []notionapi.RichText{
				{
					Type: "text",
					Text: &notionapi.Text{
						Content: c.buildTitle(),
					},
				},
			},
		},
		"Text": notionapi.RichTextProperty{
			Type: "rich_text",
			RichText: []notionapi.RichText{
				{
					Type: "text",
					Text: &notionapi.Text{
						Content: c.Text,
					},
				},
			},
		},
	}

	if c.Tags != nil {
		pageProperties["Tags"] = c.buildTagsProperties()
	}
	return &pageProperties
}
