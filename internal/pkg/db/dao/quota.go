package dao

import (
	"context"
	"notionboy/db/ent"
	"notionboy/db/ent/quota"
	"notionboy/internal/pkg/db"
	"notionboy/internal/pkg/logger"
	"strings"
	"time"
)

const (
	defaultPlan       = "free"
	defaultTotalToken = 10000
)

// QueryQuota Get Quota by user id and category
func QueryQuota(ctx context.Context, userID int) (*ent.Quota, error) {
	queryQuota := func() (*ent.Quota, error) {
		return db.GetClient().Quota.
			Query().
			Where(quota.UserIDEQ(userID)).
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
	// create a date a month later
	return db.GetClient().Quota.Create().
		SetUserID(userID).
		SetPlan(defaultPlan).
		SetToken(defaultTotalToken).
		SetResetTime(nextResetTime()).
		Exec(ctx)
}

// UpdateQuota update quota
func UpdateQuota(cli *ent.Client, ctx context.Context, userID int, tokens int64, planName string) error {
	return cli.Quota.
		Update().
		SetToken(tokens).
		SetPlan(planName).
		SetResetTime(nextResetTime()).
		Where(quota.UserIDEQ(userID)).
		Exec(ctx)
}

// IncrUsedTokenQuota increase used token quota
func IncrUsedTokenQuota(cli *ent.Client, ctx context.Context, userID int, tokens int64) error {
	return db.GetClient().Quota.
		Update().
		AddTokenUsed(tokens).
		Where(quota.UserIDEQ(userID)).
		Exec(ctx)
}

// ResetUserQuota reset user quota when the reset date is yesterday
func ResetUserQuota(ctx context.Context, resetTime time.Time) error {
	return db.GetClient().Quota.
		Update().
		SetPlan(defaultPlan).
		SetToken(defaultTotalToken).
		SetTokenUsed(0).
		SetResetTime(nextResetTime()).
		Where(quota.ResetTimeLTE(resetTime)).
		Exec(ctx)
}

func nextResetTime() time.Time {
	nextDate := time.Now().AddDate(0, 1, 0)
	return time.Date(nextDate.Year(), nextDate.Month(), nextDate.Day(), 23, 59, 59, 999, time.Local)
}
