package wxgzh

import (
	"context"
	"notionboy/config"
	"notionboy/db"
	"notionboy/notion"
	"notionboy/utils"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	log "github.com/sirupsen/logrus"
)

func messageHandler(c *gin.Context, msg *message.MixMessage) *message.Reply {

	if msg.MsgType == message.MsgType(message.EventSubscribe) {
		return replyBindingNotion(c, msg)
	}

	userID := msg.GetOpenID()
	content := transformToNotionContent(msg)
	memCache := utils.GetCache()
	userCache := memCache.Get(userID)
	log.Infof("UserID: %s, content: %s, msgType: %s, userCache: %s", userID, content, msg.MsgType, userCache)

	if msg.Content == "绑定" {
		return replyBindingNotion(c, msg)
	} else if msg.Content == "解绑" {
		return unBindingNotion(c, msg)
	}

	// 因为微信公众号没有上下文，所以使用缓存保存绑定信息
	// 如果缓存保护 userID 这个 key，说明处于绑定状态，进行绑定检测
	if memCache.Get(userID) != nil {
		token, databaseID := parseBindNotionConfig(content.Text)
		flag, msg := notion.BindingNotion(db.Account{
			NtDatabaseID: databaseID,
			NtToken:      token,
			WxUserID:     userID,
		})
		// 如果绑定成功，删除缓存的 Key
		if flag {
			memCache.Delete(userID)
		}
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(msg)}
	} else {
		// else 正常处理 Note 信息
		// 获取用户信息
		accountInfo := db.QueryAccountByWxUser(msg.GetOpenID())
		if accountInfo.ID == 0 {
			return replyBindingNotion(c, msg)
		}
		_, res := saveNoteToNotion(msg.Content, accountInfo)
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(res)}
	}
}

func transformToNotionContent(msg *message.MixMessage) *notion.Content {
	content := notion.Content{
		Text: msg.Content,
	}
	return &content
}

func replyBindingNotion(c *gin.Context, msg *message.MixMessage) *message.Reply {
	log.Warn("----- bindNotion ------")
	memCache := utils.GetCache()
	memCache.Set(msg.GetOpenID(), []string{msg.Content}, 60*time.Second)
	text := `
欢迎使用 Notion Boy，您还未绑定。
如需绑定，请在 1 分钟之内，按照下面的格式回复 Notion 的 Token 和 DatabaseID 来绑定

获取 Token 和 DatabaseID 的相关方法，请参考官方文档 https://developers.notion.com/docs/getting-started
Token 是以 "secret_" 开头的字符串，
DatabaseID 则是分享的 Page 链接的后半部分

--- 下面是具体的格式 ---
Token: secret_xxx
DatabaseID: xxxx
`
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(text)}
}

func parseBindNotionConfig(text string) (string, string) {
	log.Warn("----- parseBindNotionConfig ------")
	r := regexp.MustCompile(`Token: (?P<Token>.*) .*DatabaseID: (?P<DatabaseID>.*)`)
	res := r.FindStringSubmatch(text)
	log.Info("Parse TokenL ", res)
	return res[1], res[2]
}

func unBindingNotion(c *gin.Context, msg *message.MixMessage) *message.Reply {
	db.DeleteWxAccount(msg.GetOpenID())
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText("成功解除 Notion 绑定！")}
}

func saveNoteToNotion(msg string, accoun *db.Account) (bool, string) {
	notionContent := notion.Content{
		Text: msg,
	}
	res := notion.CreateNewRecord(context.Background(), config.Notion{BearerToken: accoun.NtToken, DatabaseID: accoun.NtDatabaseID}, notionContent)
	return strings.Contains(res, "创建 Note 成功"), res
}
