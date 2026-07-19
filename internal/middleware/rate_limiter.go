package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/web3sphere/backend/pkg/cache"
	"github.com/web3sphere/backend/pkg/response"
)

// RateLimiter limits requests per client using Redis.
func RateLimiter(redis *cache.RedisClient, maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("%s:%s", c.ClientIP(), c.FullPath())

		allowed, remaining, err := redis.CheckRateLimit(c.Request.Context(), key, maxRequests, window)
		if err != nil {
			// If Redis is down, allow the request
			c.Next()
			return
		}

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", maxRequests))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

		if !allowed {
			response.TooManyRequests(c, "Rate limit exceeded. Try again later.")
			c.Abort()
			return
		}

		c.Next()
	}
}

// StrictRateLimiter applies a stricter limit (e.g., for login, OTP endpoints).
func StrictRateLimiter(redis *cache.RedisClient) gin.HandlerFunc {
	return RateLimiter(redis, 5, 1*time.Minute)
}
