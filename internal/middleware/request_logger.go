package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/web3sphere/backend/pkg/logger"
)

// RequestLogger logs every incoming HTTP request with structured fields.
func RequestLogger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		requestID, _ := c.Get("request_id")

		fields := map[string]interface{}{
			"status":     status,
			"method":     method,
			"path":       path,
			"query":      query,
			"ip":         clientIP,
			"latency":    latency.String(),
			"latency_ms": latency.Milliseconds(),
			"user_agent": c.Request.UserAgent(),
			"request_id": requestID,
		}

		reqLogger := log.WithFields(fields)

		if status >= 500 {
			reqLogger.Errorf("Server error: %s %s", method, path)
		} else if status >= 400 {
			reqLogger.Warnf("Client error: %s %s", method, path)
		} else {
			reqLogger.Infof("%s %s", method, path)
		}
	}
}
