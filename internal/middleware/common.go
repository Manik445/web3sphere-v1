package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/web3sphere/backend/configs"
	"github.com/web3sphere/backend/pkg/response"
)

// CORS configures Cross-Origin Resource Sharing headers.
func CORS(cfg *configs.CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowedOrigins := cfg.Origins()

		allowed := false
		for _, o := range allowedOrigins {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", strings.Join(cfg.Methods(), ", "))
			c.Header("Access-Control-Allow-Headers", strings.Join(cfg.Headers(), ", "))
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", strconv.Itoa(cfg.MaxAge))
			c.Header("Access-Control-Expose-Headers", "X-Request-ID, X-RateLimit-Limit, X-RateLimit-Remaining")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// Compression enables gzip compression for responses.
func Compression() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Gin doesn't natively support compression; use gin-contrib/gzip in production.
		// This is a passthrough; the actual gzip middleware is registered via gin.Default().
		c.Next()
	}
}

// SecurityHeaders adds security-related HTTP headers.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		c.Next()
	}
}

// VersionHeader adds the API version to response headers.
func VersionHeader(version string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-API-Version", version)
		c.Next()
	}
}

// IPDetection extracts the real client IP from proxy headers.
func IPDetection() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.GetHeader("X-Forwarded-For")
		if ip == "" {
			ip = c.GetHeader("X-Real-IP")
		}
		if ip == "" {
			ip = c.ClientIP()
		}
		// Take the first IP in case of multiple proxies
		if idx := strings.Index(ip, ","); idx != -1 {
			ip = strings.TrimSpace(ip[:idx])
		}
		c.Set("client_ip", ip)
		c.Next()
	}
}

// MaintenanceMode blocks requests when maintenance mode is enabled.
func MaintenanceMode(enabled *bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if enabled != nil && *enabled {
			// Allow health check even during maintenance
			if c.Request.URL.Path == "/health" {
				c.Next()
				return
			}
			response.ServiceUnavailable(c, "System is under maintenance. Please try again later.")
			c.Abort()
			return
		}
		c.Next()
	}
}

// Timeout sets a request timeout context.
func Timeout(duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Note: Gin doesn't support per-request context timeout natively in the same way.
		// The server-level read/write timeouts handle this. This middleware adds a header.
		c.Header("X-Request-Timeout", duration.String())
		c.Next()
	}
}
