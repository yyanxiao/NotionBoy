package scheduler

import (
	"context"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/logger"
	"time"

	"github.com/go-co-op/gocron"
)

func Run() {
	s := gocron.NewScheduler(time.UTC)
	ctx := context.Background()

	// reset daily quota at 16:00 UTC (00:00 for UTC+8)every day
	if _, err := s.Every(1).Day().At("16:00").Do(dao.ResetDailyQuotaForAll, ctx); err != nil {
		logger.SugaredLogger.Errorw("scheduler ResetDailyQuotaForAll failed", "err", err)
	}

	// reset monthly quota at 16:00 UTC (00:00 for UTC+8) on the first day of every month
	if _, err := s.Every(1).MonthLastDay().At("16:02").Do(dao.ResetMonthlyQuotaForAll, ctx); err != nil {
		logger.SugaredLogger.Errorw("scheduler ResetMonthlyQuotaForAll failed", "err", err)
	}

	logger.SugaredLogger.Info("scheduler started")
	s.StartBlocking()
}
