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

func QueryAccount(ctx context.Context, userType account.UserType, userID string) (*ent.Account, error) {
	return db.GetClient().Account.
		Query().
		Where(account.UserIDEQ(userID), account.UserTypeEQ(userType)).
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
		SetNotionUserID(acc.NotionUserID).
		SetNotionUserEmail(acc.NotionUserEmail).
		OnConflict().
		UpdateNewValues().
		Exec(ctx)
}

func DeleteAccount(ctx context.Context, userType account.UserType, userID string) error {
	_, err := db.GetClient().Account.
		Delete().
		Where(account.UserIDEQ(userID), account.UserTypeEQ(userType)).
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
