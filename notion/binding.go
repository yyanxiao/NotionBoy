package notion

import (
	"context"
	"notionboy/config"
	"notionboy/db"

	"strings"

	log "github.com/sirupsen/logrus"
)

func BindingNotion(account db.Account) (bool, string) {
	log.Infof("Token: %s,\tDatabaseID: %s", account.NtToken, account.NtDatabaseID)
	if account.NtToken == "" || account.NtDatabaseID == "" {
		text := `
错误的 Token 和 DatabaseID，请按如下格式回复：
Token: secret_xxx
DatabaseID: xxxx
`
		return false, text
	} else {
		content := Content{Text: "#NotionBoy 欢迎🎉使用 Notion Boy!"}
		res := CreateNewRecord(context.Background(), config.Notion{BearerToken: account.NtToken, DatabaseID: account.NtDatabaseID}, content)
		if strings.Contains(res, "创建 Note 成功") {
			log.Debug("Token is valid, saving account.")
			db.SaveAccount(&db.Account{
				NtDatabaseID: account.NtDatabaseID,
				NtToken:      account.NtToken,
				WxUserID:     account.WxUserID,
			})
			return true, "恭喜 🎉 成功绑定 Notion！"
		} else {
			return false, "绑定 Notion 失败，无效的 Token 或 DatabaseID， 请重新绑定！"
		}
	}
}
