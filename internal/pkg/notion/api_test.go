package notion

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jomei/notionapi"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func shouldRunTest() bool {
	if os.Getenv("BEARER_TOKEN") == "" || os.Getenv("DATABASE_ID") == "" {
		return false
	}
	return true
}

func TestUpdateDatabase(t *testing.T) {
	if !shouldRunTest() {
		logrus.Info("Skip test: TestUpdateDatabase")
		return
	}

	ctx := context.Background()
	notionConfig := &NotionConfig{
		DatabaseID:  os.Getenv("DATABASE_ID"),
		BearerToken: os.Getenv("BEARER_TOKEN"),
	}

	respMsg, err := updateDatabase(ctx, notionConfig, &notionapi.DatabaseUpdateRequest{})
	assert.Nil(t, err, respMsg)
}

func TestCreateNewRecord(t *testing.T) {
	if !shouldRunTest() {
		logrus.Info("Skip test: TestCreateNewRecord")
		return
	}

	ctx := context.Background()
	notionConfig := &NotionConfig{
		DatabaseID:  os.Getenv("DATABASE_ID"),
		BearerToken: os.Getenv("BEARER_TOKEN"),
	}
	content := &Content{
		Tags: []string{"test"},
		Text: "This is test",
	}

	respMsg, err := CreateNewRecord(ctx, notionConfig, content)
	assert.Nil(t, err, respMsg)
}

func TestParseContent(t *testing.T) {
	c := Content{
		Text: "#hello #world\n #我\t#🤔 #end",
	}
	c.parseTags(context.TODO())
	assert.Equal(t, []string{"hello", "world", "我", "🤔", "end"}, c.Tags)
}

func TestParseFulltextContent(t *testing.T) {
	url1 := "https://url1.com/T02RH5Q0K/DJ9TDT8KV"
	url2 := "http://url2.234/T02RH5Q0K"
	url3 := "url3.abc/T02RH5Q0K/DJ9TDT8KV"
	c := Content{
		Text: fmt.Sprintf("#全文 %s\t%s\n%s", url1, url2, url3),
	}
	c.parseTags(context.TODO())
	assert.Equal(t, url1, c.Fulltext.URL)
}
