package middleware

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Format-C-eft/utils/logger"
)

func Recovery() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(ginContext.Request, false)
				if brokenPipe {
					logger.ErrorKV(
						ginContext.Request.Context(),
						ginContext.Request.URL.Path,
						"error", err,
						"request", string(httpRequest),
					)

					// If the connection is dead, we can't write a status to it.
					ginContext.Error(err.(error)) // nolint: errcheck
					ginContext.Abort()
					return
				}

				logger.ErrorKV(
					ginContext.Request.Context(),
					"[Recovery from panic]",
					"error", err,
					"request", string(httpRequest),
					"stack", string(debug.Stack()),
				)
				ginContext.AbortWithStatus(http.StatusInternalServerError)
			}
		}()

		ginContext.Next()
	}
}
