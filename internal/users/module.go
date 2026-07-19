package users

import (
	"github.com/gin-gonic/gin"
	"github.com/web3sphere/backend/pkg/cache"
	"github.com/web3sphere/backend/pkg/jwt"
	"github.com/web3sphere/backend/pkg/logger"
	"gorm.io/gorm"
)

// Setup wires up the Users module dependencies and registers its routes.
func Setup(
	router *gin.RouterGroup,
	db *gorm.DB,
	jwtManager *jwt.Manager,
	redisClient *cache.RedisClient,
	log *logger.Logger,
) {
	repo := NewRepository(db)
	svc := NewService(repo, log)
	ctrl := NewController(svc)

	RegisterRoutes(router, ctrl, jwtManager, redisClient)
}
