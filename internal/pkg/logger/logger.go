package logger

import (
	"notionboy/internal/pkg/config"

	"go.uber.org/zap/zapcore"

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
	level, _ := zapcore.ParseLevel(config.GetConfig().Log.Level)
	if level == zapcore.DebugLevel {
		Logger, _ = zap.NewDevelopment()
	} else {
		Logger, _ = zap.NewProduction()
	}

	SugaredLogger = Logger.Sugar()
}
