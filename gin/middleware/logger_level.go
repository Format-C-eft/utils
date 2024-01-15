package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/Format-C-eft/utils/headers"
	"github.com/Format-C-eft/utils/logger"
)

func ChangeLoggerLevel() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		zapLevel, ok := logger.StrToZapLevel(ginContext.GetHeader(headers.XLogLevel))
		if !ok {
			return
		}

		ctx := ginContext.Request.Context()

		newLogger := logger.CloneWithLevel(ctx, zapLevel)
		ginContext.Request = ginContext.Request.Clone(
			logger.AttachLogger(ctx, newLogger),
		)

		ginContext.Next()
	}
}
