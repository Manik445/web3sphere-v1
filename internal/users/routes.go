package users

import (
	"github.com/gin-gonic/gin"
	"github.com/web3sphere/backend/internal/middleware"
	"github.com/web3sphere/backend/pkg/cache"
	"github.com/web3sphere/backend/pkg/jwt"
)

// RegisterRoutes registers the user routes.
func RegisterRoutes(router *gin.RouterGroup, ctrl *Controller, jwtManager *jwt.Manager, redis *cache.RedisClient) {
	usersRoutes := router.Group("/users")
	usersRoutes.Use(middleware.Auth(jwtManager, redis))
	{
		usersRoutes.GET("/:id", ctrl.GetProfile)
		usersRoutes.PUT("/me", ctrl.UpdateProfile)
	}
}
