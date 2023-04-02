package scheduler

import (
	"context"
	"notionboy/internal/pkg/db/dao"
	"time"
)

func resetUserQuota(ctx context.Context) error {
	// reset user quota when the reset date is yesterday
	return dao.ResetUserQuota(ctx, time.Now())
}
