package logger

import (
	"context"
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKey struct{}

var (
	// global глобальный экземпляр логгера.
	global            *zap.SugaredLogger
	defaultLevel      = zap.NewAtomicLevelAt(zap.ErrorLevel)
	attachedLoggerKey = &ctxKey{}
)

func init() {
	SetLogger(New(defaultLevel))
}

func New(level zapcore.LevelEnabler, options ...zap.Option) *zap.SugaredLogger {
	return NewWithSink(level, os.Stdout, options...)
}

func NewWithSink(level zapcore.LevelEnabler, sink io.Writer, options ...zap.Option) *zap.SugaredLogger {
	if level == nil {
		level = defaultLevel
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}),
		zapcore.AddSync(sink),
		level,
	)

	return zap.New(core, options...).Sugar()
}

func SetLevel(l zapcore.Level) {
	defaultLevel.SetLevel(l)
}

func SetLogger(l *zap.SugaredLogger) {
	global = l
}

func GetLogger() *zap.SugaredLogger {
	return global
}

func Debug(ctx context.Context, args ...interface{}) {
	fromContext(ctx).Debug(args...)
}

func DebugF(ctx context.Context, format string, args ...interface{}) {
	fromContext(ctx).Debugf(format, args...)
}

func DebugKV(ctx context.Context, message string, kvs ...interface{}) {
	fromContext(ctx).Debugw(message, fixError(kvs...)...)
}

func Info(ctx context.Context, args ...interface{}) {
	fromContext(ctx).Info(args...)
}

func InfoF(ctx context.Context, format string, args ...interface{}) {
	fromContext(ctx).Infof(format, args...)
}

func InfoKV(ctx context.Context, message string, kvs ...interface{}) {
	fromContext(ctx).Infow(message, fixError(kvs...)...)
}

func Warn(ctx context.Context, args ...interface{}) {
	fromContext(ctx).Warn(args...)
}

func WarnF(ctx context.Context, format string, args ...interface{}) {
	fromContext(ctx).Warnf(format, args...)
}

func WarnKV(ctx context.Context, message string, kvs ...interface{}) {
	fromContext(ctx).Warnw(message, fixError(kvs...)...)
}

func Error(ctx context.Context, args ...interface{}) {
	fromContext(ctx).Error(args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	fromContext(ctx).Errorf(format, args...)
}

func ErrorKV(ctx context.Context, message string, kvs ...interface{}) {
	fromContext(ctx).Errorw(message, fixError(kvs...)...)
}

func Fatal(ctx context.Context, args ...interface{}) {
	fromContext(ctx).Fatal(args...)
}

func FatalF(ctx context.Context, format string, args ...interface{}) {
	fromContext(ctx).Fatalf(format, args...)
}

func FatalKV(ctx context.Context, message string, kvs ...interface{}) {
	fromContext(ctx).Fatalw(message, kvs...)
}

func Panic(ctx context.Context, args ...interface{}) {
	fromContext(ctx).Panic(args...)
}
func PanicF(ctx context.Context, format string, args ...interface{}) {
	fromContext(ctx).Panicf(format, args...)
}

func PanicKV(ctx context.Context, message string, kvs ...interface{}) {
	fromContext(ctx).Panicw(message, kvs...)
}

func fixError(args ...interface{}) []interface{} {
	if args == nil {
		return nil
	}

	result := make([]interface{}, 0, len(args))
	for _, arg := range args {
		if val, ok := arg.(error); ok {
			result = append(result, val.Error())
			continue
		}
		result = append(result, arg)
	}

	return result
}
