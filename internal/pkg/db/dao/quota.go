package dao

import (
	"context"
	"notionboy/db/ent"
	"notionboy/db/ent/quota"
	"notionboy/internal/pkg/db"
	"notionboy/internal/pkg/logger"
	"strings"
)

// QueryQuota Get Quota by user id and category
func QueryQuota(ctx context.Context, userID int, category quota.Category) (*ent.Quota, error) {
	queryQuota := func() (*ent.Quota, error) {
		return db.GetClient().Quota.
			Query().
			Where(quota.UserIDEQ(userID), quota.CategoryEQ(category)).
			Only(ctx)
	}
	q, err := queryQuota()
	if err != nil {
		// if not found, create a new one
		if strings.Contains(err.Error(), "not found") {
			if err = initQuota(ctx, userID); err != nil {
				return nil, err
			}
			q, err = queryQuota()
		} else {
			logger.SugaredLogger.Errorw("query quota failed", "err", err)
			return nil, err
		}
	}

	return q, err
}

// initQuota init quota with default value
func initQuota(ctx context.Context, userID int) error {
	quotas := make([]*ent.QuotaCreate, 0)
	newQuota := func(category quota.Category, daily int) *ent.QuotaCreate {
		q := db.GetClient().Quota.Create()
		q.SetUserID(userID)
		q.SetCategory(category)
		q.SetDaily(daily)
		return q
	}
	quotas = append(quotas, newQuota(quota.CategoryChatgpt, 10))

	err := db.GetClient().Quota.
		CreateBulk(quotas...).
		Exec(ctx)
	if err != nil {
		logger.SugaredLogger.Errorw("init quota failed", "err", err)
	}
	return err
}

// SaveQuota Save Quota
func SaveQuota(ctx context.Context, q *ent.Quota) error {
	return db.GetClient().Quota.
		Create().
		SetUserID(q.UserID).
		SetDaily(q.Daily).
		SetMonthly(q.Monthly).
		SetYearly(q.Yearly).
		SetDailyUsed(q.DailyUsed).
		SetMonthlyUsed(q.MonthlyUsed).
		SetYearlyUsed(q.YearlyUsed).
		OnConflict().
		UpdateNewValues().
		Exec(ctx)
}

// IncrDailyQuota increment daily quota
func IncrDailyQuota(ctx context.Context, userID int, category quota.Category) error {
	dailyUsed, err := GetDailyUsedQuota(ctx, userID, category)
	if err != nil {
		return err
	}
	return db.GetClient().Quota.
		Update().
		SetDailyUsed(dailyUsed+1).
		Where(quota.UserIDEQ(userID), quota.CategoryEQ(category)).
		Exec(ctx)
}

// ResetDailyQuota reset daily quota
func ResetDailyQuota(ctx context.Context, userID int, category quota.Category) error {
	return db.GetClient().Quota.
		Update().
		SetDailyUsed(0).
		Where(quota.UserIDEQ(userID), quota.CategoryEQ(category)).
		Exec(ctx)
}

// GetDailyUsedQuota get daily quota
func GetDailyUsedQuota(ctx context.Context, userID int, category quota.Category) (int, error) {
	q, err := QueryQuota(ctx, userID, category)
	if err != nil {
		return 0, err
	}
	return q.DailyUsed, nil
}
