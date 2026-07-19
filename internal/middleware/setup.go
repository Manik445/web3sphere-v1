package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/web3sphere/backend/configs"
	"github.com/web3sphere/backend/pkg/logger"
)

// ApplyGlobal applies all global middlewares to the Gin engine in the correct order.
func ApplyGlobal(router *gin.Engine, cfg *configs.Config, log *logger.Logger) {
	router.Use(RequestID())
	router.Use(RequestLogger(log))
	router.Use(Recovery(log))
	router.Use(CORS(&cfg.CORS))
	router.Use(Compression())
	router.Use(IPDetection())
	router.Use(SecurityHeaders())
	router.Use(VersionHeader(cfg.App.Version))
	router.Use(Timeout(cfg.Server.ReadTimeout))
}
