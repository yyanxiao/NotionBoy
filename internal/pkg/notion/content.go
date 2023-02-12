package notion

import (
	"context"
	"notionboy/db/ent"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"
	"regexp"
	"strings"
	"time"

	"github.com/jomei/notionapi"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Content struct {
	Tags          []string         `json:"tags"`
	Text          string           `json:"text"`
	NotionPageID  string           `json:"notion_page_id"`
	IsMedia       bool             `json:"is_media"`
	Media         MediaContent     `json:"media"`
	IsChatContent bool             `json:"is_chat_content"`
	ChatContent   ChatContent      `json:"chat_content"`
	Medias        []*MediaContent  `json:"medias"`
	Account       *ent.Account     `json:"account"`
	Zlib          *ZlibContent     `json:"zlib"`
	FullText      *FulltextContent `json:"fulltext"`
}

// Process 从 text 提取 tags，配置全文
func (c *Content) Process(ctx context.Context) {
	r, _ := regexp.Compile(`#(.+?)($|\s)`)
	match := r.FindAllStringSubmatch(c.Text, -1)
	if len(match) > 0 {
		tags := make([]string, 0)
		for _, m := range match {
			tag := strings.Trim(m[1], "# ")
			tags = append(tags, tag)
			if strings.ToUpper(tag) == config.CMD_FULLTEXT {
				c.parseFulltext(ctx)
			}
		}
		c.Tags = tags
	}
}

func parseUrlFromText(text string) string {
	r, _ := regexp.Compile(`https?://(www.)?[-a-zA-Z0-9@:%._+~#=]{2,256}.[a-z]{2,4}\b([-a-zA-Z0-9@:%_+.~#?&//=]*)`)
	match := r.FindAllStringSubmatch(text, -1)
	if len(match) > 0 {
		// only save last url
		for _, m := range match {
			return m[0]
		}
	}
	return ""
}

func (c *Content) parseFulltext(ctx context.Context) {
	url := parseUrlFromText(c.Text)
	if url != "" {
		c.FullText = &FulltextContent{
			URL:          url,
			NotionPageID: c.NotionPageID,
			Account:      c.Account,
		}
		c.FullText.ProcessFulltext(ctx)
	}
}

func (c *Content) buildTitle() string {
	title := c.Text
	if c.IsMedia && c.Text == "" {
		loc, _ := time.LoadLocation("Asia/Shanghai")
		title = cases.Title(language.English).String(c.Media.Type) + " " + time.Now().UTC().In(loc).Format(time.RFC3339)
	}
	if c.FullText != nil && c.FullText.Title != "" {
		title = c.FullText.Title
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
	if c.Account != nil &&
		c.Account.NotionUserID != "" &&
		c.Account.DatabaseID != config.GetConfig().NotionTestPage.DatabaseID {
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
								ID: notionapi.UserID(c.Account.NotionUserID),
							},
						},
					},
				},
			},
		})
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
	if c.Zlib != nil {
		blocks = append(blocks, c.Zlib.BuildBlocks()...)
	}
	if c.FullText != nil {
		blocks = append(blocks, c.FullText.BuildBlocks()...)
	}
	return blocks
}

func (c *Content) BuildPageProperties() *notionapi.Properties {
	setRichText := func(text string) notionapi.RichTextProperty {
		return notionapi.RichTextProperty{
			Type: "rich_text",
			RichText: []notionapi.RichText{
				{
					Type: "text",
					Text: &notionapi.Text{
						Content: text,
					},
				},
			},
		}
	}
	// if there is no text, do no set title and text
	pageProperties := notionapi.Properties{}
	if c.Text != "" {
		pageProperties["Name"] = notionapi.TitleProperty{
			Type: "title",
			Title: []notionapi.RichText{
				{
					Type: "text",
					Text: &notionapi.Text{
						Content: c.buildTitle(),
					},
				},
			},
		}
		pageProperties["Text"] = setRichText(c.Text)
	}

	if c.FullText != nil {
		pageProperties["Author"] = setRichText(c.FullText.Author)
		pageProperties["Summary"] = setRichText(c.FullText.Summary)
		pageProperties["URL"] = notionapi.URLProperty{
			Type: "url",
			URL:  c.FullText.URL,
		}

		var d notionapi.Date
		if err := d.UnmarshalText([]byte(c.FullText.PublishDate)); err == nil {
			pageProperties["PublishDate"] = notionapi.DateProperty{
				Type: "date",
				Date: &notionapi.DateObject{
					Start: &d,
				},
			}
		} else {
			logger.SugaredLogger.Errorw("parse publish date error", "err", err)
		}
	}

	if c.Tags != nil {
		pageProperties["Tags"] = c.buildTagsProperties()
	}
	return &pageProperties
}
