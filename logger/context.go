package logger

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// fromContext достает логгер из контекста. Если в контексте логгер не
// обнаруживается - возвращает глобальный логгер.
func fromContext(ctx context.Context) *zap.SugaredLogger {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ctx = ginCtx.Request.Context()
	}

	if logger, ok := ctx.Value(attachedLoggerKey).(*zap.SugaredLogger); ok {
		return logger
	}

	return global
}
