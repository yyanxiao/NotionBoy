package dao

import (
	"context"
	"notionboy/db/ent"
	"notionboy/db/ent/account"
	"notionboy/internal/pkg/db"
)

func QueryAccountByWxUser(ctx context.Context, wxUserID string) (*ent.Account, error) {
	return db.GetClient().Account.
		Query().
		Where(account.UserIDEQ(wxUserID), account.UserTypeEQ(account.UserTypeWechat)).
		Only(ctx)
}

func SaveAccount(ctx context.Context, acc *ent.Account) error {
	return db.GetClient().Account.
		Create().
		SetUserID(acc.UserID).
		SetUserType(acc.UserType).
		SetDatabaseID(acc.DatabaseID).
		SetAccessToken(acc.AccessToken).
		SetIsLatestSchema(acc.IsLatestSchema).
		OnConflict().
		UpdateNewValues().
		Exec(ctx)
}

func DeleteWxAccount(ctx context.Context, wxUserID string) error {
	_, err := db.GetClient().Account.
		Delete().
		Where(account.UserIDEQ(wxUserID), account.UserTypeEQ(account.UserTypeWechat)).
		Exec(ctx)
	return err
}

func UpdateIsLatestSchema(ctx context.Context, databaseID string, isLatest bool) error {
	return db.GetClient().Account.
		Update().
		SetIsLatestSchema(isLatest).
		Where(account.DatabaseIDEQ(databaseID)).
		Exec(ctx)
}
