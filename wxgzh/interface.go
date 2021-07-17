package wxgzh

import "notionboy/db"

type Wxgzh interface {
	Run()
	bindingNotion(token, databaseID, userID string) (bool, string)
	unBindingNotion(userID string) (bool, string)
	saveNoteToNotion(msg string, accoun *db.Account) (bool, string)
}
