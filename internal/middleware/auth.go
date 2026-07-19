package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/web3sphere/backend/pkg/cache"
	pkgjwt "github.com/web3sphere/backend/pkg/jwt"
	"github.com/web3sphere/backend/pkg/response"
)

// Auth validates the JWT access token and sets user info in context.
func Auth(jwtManager *pkgjwt.Manager, redis *cache.RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Unauthorized(c, "Invalid authorization format. Use: Bearer <token>")
			c.Abort()
			return
		}

		claims, err := jwtManager.ValidateAccessToken(parts[1])
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Check if token is blacklisted
		blacklisted, err := redis.IsTokenBlacklisted(c.Request.Context(), claims.ID)
		if err != nil || blacklisted {
			response.Unauthorized(c, "Token has been revoked")
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("token_id", claims.ID)

		c.Next()
	}
}

// OptionalAuth tries to authenticate but doesn't fail if no token is provided.
func OptionalAuth(jwtManager *pkgjwt.Manager, redis *cache.RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			if claims, err := jwtManager.ValidateAccessToken(parts[1]); err == nil {
				blacklisted, _ := redis.IsTokenBlacklisted(c.Request.Context(), claims.ID)
				if !blacklisted {
					c.Set("user_id", claims.UserID)
					c.Set("email", claims.Email)
					c.Set("role", claims.Role)
					c.Set("token_id", claims.ID)
				}
			}
		}

		c.Next()
	}
}

// GetUserID extracts the user ID from gin context (set by Auth middleware).
func GetUserID(c *gin.Context) string {
	if v, exists := c.Get("user_id"); exists {
		return v.(string)
	}
	return ""
}

// GetUserRole extracts the user role from gin context.
func GetUserRole(c *gin.Context) string {
	if v, exists := c.Get("role"); exists {
		return v.(string)
	}
	return ""
}

// RequireVerifiedEmail ensures the user's email is verified.
func RequireVerifiedEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		verified, exists := c.Get("email_verified")
		if !exists || !verified.(bool) {
			response.Error(c, http.StatusForbidden, "Email verification required", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
