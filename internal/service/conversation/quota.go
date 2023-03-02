package conversation

import (
	"context"
	"notionboy/db/ent"
	"notionboy/db/ent/quota"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/logger"
)

func checkRateLimit(acc *ent.Account, qt *ent.Quota) bool {
	if qt.DailyUsed >= qt.Daily || qt.MonthlyUsed >= qt.Monthly {
		logger.SugaredLogger.Debugw("Hit rate limit", "account", acc.ID, "daily_used", qt.DailyUsed, "daily", qt.Daily)
		return true
	}
	logger.SugaredLogger.Debugw("Not hit rate limit", "account", acc.ID, "daily_used", qt.DailyUsed, "daily", qt.Daily, "category", qt.Category)

	return false
}

func loadQuota(ctx context.Context, acc *ent.Account) (*ent.Quota, error) {
	return dao.QueryQuota(ctx, acc.ID, quota.CategoryChatgpt)
}
