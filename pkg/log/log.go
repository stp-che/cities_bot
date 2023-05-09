package log

import (
	"context"

	"go.uber.org/zap"
)

var logger *zap.Logger

func getLogger() *zap.Logger {
	if logger == nil {
		var err error
		logger, err = zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
	}

	return logger
}

func Info(_ context.Context, msg string) {
	getLogger().Info(msg)
}

func Warn(_ context.Context, msg string) {
	getLogger().Warn(msg)
}

func Fatal(_ context.Context, msg string) {
	getLogger().Fatal(msg)
}
