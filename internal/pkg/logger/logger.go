package logger

import (
	"notionboy/internal/pkg/config"

	"go.uber.org/zap"
)

var (
	SugaredLogger *zap.SugaredLogger
	Logger        *zap.Logger
)

const (
	LevelDebug = "debug"
)

func init() {
	if config.GetConfig().Log.Level == LevelDebug {
		Logger, _ = zap.NewDevelopment()
	} else {
		Logger, _ = zap.NewProduction()
	}
	SugaredLogger = Logger.Sugar()
}
