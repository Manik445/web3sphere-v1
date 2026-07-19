package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/web3sphere/backend/internal/middleware"
	"github.com/web3sphere/backend/pkg/cache"
	"github.com/web3sphere/backend/pkg/jwt"
)

// RegisterRoutes registers the authentication routes.
func RegisterRoutes(router *gin.RouterGroup, ctrl *Controller, jwtManager *jwt.Manager, redis *cache.RedisClient) {
	authRoutes := router.Group("/auth")
	{
		// Rate limited endpoints
		strictLimit := middleware.StrictRateLimiter(redis)
		
		authRoutes.POST("/signup", strictLimit, ctrl.Signup)
		authRoutes.POST("/login", strictLimit, ctrl.Login)
		authRoutes.POST("/verify-email", strictLimit, ctrl.VerifyEmail)
		authRoutes.POST("/resend-otp", strictLimit, ctrl.ResendOTP)
		authRoutes.POST("/forgot-password", strictLimit, ctrl.ForgotPassword)
		authRoutes.POST("/reset-password", strictLimit, ctrl.ResetPassword)
		
		// Standard endpoints
		authRoutes.POST("/refresh", ctrl.RefreshToken)
		
		// Protected endpoints
		protected := authRoutes.Group("")
		protected.Use(middleware.Auth(jwtManager, redis))
		{
			protected.POST("/logout", ctrl.Logout)
			protected.GET("/me", ctrl.Me)
		}
	}
}
