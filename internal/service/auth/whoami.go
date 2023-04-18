package auth

import (
	"context"
	"fmt"
	"time"

	"notionboy/db/ent"
	"notionboy/db/ent/account"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/logger"

	"github.com/google/uuid"
)

type MyInfo struct {
	Account *ent.Account
	Quota   *ent.Quota
}

type MyInfoDTO struct {
	UserID           uuid.UUID `json:"user_id"`            // uuid of user
	UserType         string    `json:"user_type"`          // type of user, "wechat" or "telegram"
	NotionDatabaseID string    `json:"notion_database_id"` // id of notion database, if not exist, it will be ""
	Plan             string    `json:"plan"`               // plan of user, default to "free"
	TotalToken       int64     `json:"toatl_token"`        // quota of user
	UsedToken        int64     `json:"used_token"`         // used quota of user, default to 0
	ResetTime        time.Time `json:"reset_time"`         // reset time of quota
	ApiKey           string    `json:"api_key"`            // api key of user
}

func (m *MyInfo) ToDTO() *MyInfoDTO {
	return &MyInfoDTO{
		UserID:           m.Account.UUID,
		UserType:         m.Account.UserType.String(),
		NotionDatabaseID: m.Account.DatabaseID,
		Plan:             m.Quota.Plan,
		TotalToken:       m.Quota.Token,
		UsedToken:        m.Quota.TokenUsed,
		ResetTime:        m.Quota.ResetTime,
		ApiKey:           m.Account.APIKey.String(),
	}
}

func (m *MyInfoDTO) String() string {
	// concatentate all fields and split by \n
	return fmt.Sprintf("用户ID: %s\n\n用户类型: %s\n\nNotionDB: %s\n\n订阅计划: %s\n\nToken总量: %d\n\n已用Token: %d\n\n额度重置时间: %s\n\nAPIKey: %s\n",
		m.UserID, m.UserType, m.NotionDatabaseID, m.Plan, m.TotalToken, m.UsedToken, m.ResetTime.Format(time.RFC3339), m.ApiKey)
}

func WhoAmI(ctx context.Context, userType account.UserType, userId string) (*MyInfoDTO, error) {
	acc, err := dao.QueryAccount(ctx, userType, userId)
	if err != nil {
		logger.SugaredLogger.Errorw("query account failed", "error", err, "userType", userType, "userId", userId)
		return nil, err
	}

	// if no api key, generate one
	if acc.APIKey == uuid.Nil {
		apiKey := uuid.New()
		err = dao.UpdateAccountApiKey(ctx, acc.UUID, apiKey)
		if err != nil {
			logger.SugaredLogger.Errorw("Update account set api key failed", "error", err, "account", acc)
			return nil, err
		}
		acc.APIKey = apiKey
	}

	quota, err := dao.QueryQuota(ctx, acc.ID)
	if err != nil {
		logger.SugaredLogger.Errorw("query quota failed", "error", err, "account", acc)
		return nil, err
	}

	myInfo := &MyInfo{
		Account: acc,
		Quota:   quota,
	}

	return myInfo.ToDTO(), nil
}
