package handler

import (
	"context"
	"fmt"
	"notionboy/db/ent"
	"notionboy/db/ent/account"
	"notionboy/internal/pkg/db/dao"
	"strconv"

	tele "gopkg.in/telebot.v3"
)

func queryUserAccount(ctx context.Context, c tele.Context) (*ent.Account, error) {
	sender := c.Sender()
	if sender == nil {
		return nil, fmt.Errorf("User do not exist")
	}
	return dao.QueryAccount(ctx, account.UserTypeTelegram, strconv.FormatInt(sender.ID, 10))
}
