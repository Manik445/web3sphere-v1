package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/web3sphere/backend/pkg/logger"
	"github.com/web3sphere/backend/pkg/response"
)

// Recovery recovers from panics and logs the stack trace.
func Recovery(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				requestID, _ := c.Get("request_id")

				log.WithFields(map[string]interface{}{
					"error":      err,
					"stack":      stack,
					"request_id": requestID,
					"path":       c.Request.URL.Path,
					"method":     c.Request.Method,
				}).Errorf("Panic recovered: %v", err)

				response.Error(c, http.StatusInternalServerError, "Internal server error", nil)
				c.Abort()
			}
		}()
		c.Next()
	}
}
