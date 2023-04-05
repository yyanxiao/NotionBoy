package scheduler

import (
	"context"
	"time"

	"notionboy/internal/pkg/db/dao"
)

func resetUserQuota(ctx context.Context) error {
	// reset user quota when the reset date is yesterday
	return dao.ResetUserQuota(ctx, time.Now())
}
