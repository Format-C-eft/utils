package logger

import (
	"context"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// StrToZapLevel - convert string to zap.level
func StrToZapLevel(str string) (zapcore.Level, bool) {
	switch strings.ToLower(str) {
	case "debug":
		return zapcore.DebugLevel, true
	case "info":
		return zapcore.InfoLevel, true
	case "warn":
		return zapcore.WarnLevel, true
	case "error":
		return zapcore.ErrorLevel, true
	}

	return zapcore.WarnLevel, false
}

// CloneWithLevel - clone logger with level
func CloneWithLevel(ctx context.Context, newLevel zapcore.Level) *zap.SugaredLogger {
	return fromContext(ctx).
		Desugar().
		WithOptions(
			zap.WrapCore(
				func(core zapcore.Core) zapcore.Core {
					return &coreWithLevel{
						core,
						newLevel,
					}
				},
			),
		).
		Sugar()
}

// AttachLogger - context attach logger
func AttachLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, attachedLoggerKey, logger)
}
