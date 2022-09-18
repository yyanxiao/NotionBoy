package notion

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {
	c := &Content{
		Text: "#hello #world\n #我\t#🤔 #end",
	}
	c.Process(context.TODO())
	assert.Equal(t, []string{"hello", "world", "我", "🤔", "end"}, c.Tags)
}

// func TestParseFulltextContent(t *testing.T) {
// 	url1 := "https://url1.com/T02RH5Q0K/DJ9TDT8KV"
// 	url2 := "http://url2.234/T02RH5Q0K"
// 	url3 := "url3.abc/T02RH5Q0K/DJ9TDT8KV"
// 	c := Content{
// 		Text: fmt.Sprintf("#全文 %s\t%s\n%s", url1, url2, url3),
// 	}
// 	c.parseFulltextURL(context.TODO(), "全文")
// 	assert.Equal(t, url1, c.Fulltext.URL)
// }
