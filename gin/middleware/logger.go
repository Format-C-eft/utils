package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Format-C-eft/utils/logger"
)

func Logger() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		start := time.Now()

		logger.DebugKV(
			ginContext,
			"request info",
			"method", ginContext.Request.Method,
			"path", ginContext.Request.URL.Path,
			"query", ginContext.Request.URL.RawQuery,
			"header", ginContext.Request.Header,
		)

		ginContext.Next()

		if len(ginContext.Errors) > 0 {
			// you can add fields if required
			for _, err := range ginContext.Errors.Errors() {
				logger.ErrorKV(ginContext, "errors of context", "err", err)
			}
		} else {
			logger.DebugKV(
				ginContext,
				"response info",
				"status", ginContext.Writer.Status(),
				"header", ginContext.Writer.Header(),
				"latency", time.Now().Sub(start),
			)
		}
	}
}
