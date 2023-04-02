package scheduler

import (
	"context"
	"notionboy/internal/pkg/logger"
	"time"

	"github.com/go-co-op/gocron"
)

func Run() {
	s := gocron.NewScheduler(time.UTC)
	ctx := context.Background()

	// reset user quota every day at 00:00
	if _, err := s.Day().Every(1).At("00:00").Do(resetUserQuota, ctx); err != nil {
		logger.SugaredLogger.Errorw("scheduler ResetMonthlyQuotaForAll failed", "err", err)
	}
	logger.SugaredLogger.Info("scheduler started")
	s.StartBlocking()
}
