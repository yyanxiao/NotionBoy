package dao

import (
	"context"

	"notionboy/db/ent"
	"notionboy/internal/pkg/db"
)

const dummyUserID = "dummy_user_id"

func SaveSession(ctx context.Context, session []byte) error {
	return db.GetClient().WechatSession.
		Create().
		SetDummyUserID(dummyUserID).
		SetSession(session).
		OnConflict().
		UpdateNewValues().
		Exec(ctx)
}

func GetSession(ctx context.Context) (*ent.WechatSession, error) {
	return db.GetClient().WechatSession.Query().First(ctx)
}
