package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/web3sphere/backend/configs"
	"github.com/web3sphere/backend/pkg/cache"
	"github.com/web3sphere/backend/pkg/jwt"
	"github.com/web3sphere/backend/pkg/logger"
	"github.com/web3sphere/backend/pkg/mailer"
	"gorm.io/gorm"
)

// Setup wires up the Auth module dependencies and registers its routes.
func Setup(
	router *gin.RouterGroup,
	db *gorm.DB,
	cfg *configs.Config,
	redisClient *cache.RedisClient,
	jwtManager *jwt.Manager,
	mailSvc mailer.Mailer,
	log *logger.Logger,
) {
	repo := NewRepository(db)
	svc := NewService(repo, cfg, redisClient, jwtManager, mailSvc, log)
	ctrl := NewController(svc)
	
	RegisterRoutes(router, ctrl, jwtManager, redisClient)
}
