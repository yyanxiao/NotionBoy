package dao

import (
	"context"
	"notionboy/db/ent"
	"notionboy/db/ent/account"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/db"
	"notionboy/internal/pkg/logger"
	"strings"
)

// QueryAccountByWxUser Get Account by wx user id
func QueryAccountByWxUser(ctx context.Context, wxUserID string) (*ent.Account, error) {
	return QueryAccount(ctx, account.UserTypeWechat, wxUserID)
}

// QueryAccount Get Account by user id and user type
func QueryAccount(ctx context.Context, userType account.UserType, userID string) (*ent.Account, error) {
	query := func() (*ent.Account, error) {
		return db.GetClient().Account.
			Query().
			Where(account.UserIDEQ(userID), account.UserTypeEQ(userType)).
			Only(ctx)
	}
	acc, err := query()
	if err != nil {
		// if not found, create a new one
		if strings.Contains(err.Error(), "not found") {
			err = initAccount(ctx, userID, userType)
			if err != nil {
				return nil, err
			}
			acc, err = query()
		} else {
			return nil, err
		}
	}

	// if the account do not have a database id, set it to test database id
	if acc.DatabaseID == "" {
		acc.DatabaseID = config.GetConfig().NotionTestPage.DatabaseID
		acc.AccessToken = config.GetConfig().NotionTestPage.Token
	}

	return acc, err
}

func initAccount(ctx context.Context, userID string, userType account.UserType) error {
	err := db.GetClient().Account.
		Create().
		SetUserID(userID).
		SetUserType(userType).
		OnConflict().
		UpdateNewValues().
		Exec(ctx)
	if err != nil {
		logger.SugaredLogger.Errorw("init account failed", "err", err)
	}
	return err
}

// SaveAccount Save Account
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

// SaveBasicAccount Save Basic Account info
func SaveBasicAccount(ctx context.Context, userType account.UserType, userID string) error {
	return db.GetClient().Account.
		Create().
		SetUserID(userID).
		SetUserType(userType).
		SetIsLatestSchema(true).
		OnConflict().
		UpdateNewValues().
		Exec(ctx)
}

// DeleteAccount Delete Account
func DeleteAccount(ctx context.Context, userType account.UserType, userID string) error {
	_, err := db.GetClient().Account.
		Delete().
		Where(account.UserIDEQ(userID), account.UserTypeEQ(userType)).
		Exec(ctx)
	return err
}

// ClearNotionAuthInfo Clear Notion Auth Info
func ClearNotionAuthInfo(ctx context.Context, userType account.UserType, userID string) error {
	err := db.GetClient().Account.
		Update().
		SetDatabaseID("").
		SetAccessToken("").
		SetNotionUserID("").
		SetNotionUserEmail("").
		Where(account.UserIDEQ(userID), account.UserTypeEQ(userType)).
		Exec(ctx)
	return err
}

// UpdateIsLatestSchema update is latest schema
func UpdateIsLatestSchema(ctx context.Context, databaseID string, isLatest bool) error {
	return db.GetClient().Account.
		Update().
		SetIsLatestSchema(isLatest).
		Where(account.DatabaseIDEQ(databaseID)).
		Exec(ctx)
}
