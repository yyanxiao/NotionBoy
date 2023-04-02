package auth

import (
	"context"
	"notionboy/db/ent/account"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/logger"

	"github.com/google/uuid"
)

func GenerateApiKey(ctx context.Context, userType account.UserType, userId string) (string, error) {
	acc, err := dao.QueryAccount(ctx, userType, userId)
	if err != nil {
		logger.SugaredLogger.Errorw("GenerateApiKey query account failed", "error", err, "userType", userType, "userId", userId)
		return "", err
	}

	apiKey := uuid.New()
	if err := dao.UpdateAccountApiKey(ctx, acc.UUID, apiKey); err != nil {
		return "", err
	}
	return apiKey.String(), nil
}
