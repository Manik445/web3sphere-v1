package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/web3sphere/backend/configs"
	"github.com/web3sphere/backend/internal/auth"
	"github.com/web3sphere/backend/internal/middleware"
	"github.com/web3sphere/backend/internal/users"
	"github.com/web3sphere/backend/pkg/cache"
	"github.com/web3sphere/backend/pkg/database"
	"github.com/web3sphere/backend/pkg/jwt"
	"github.com/web3sphere/backend/pkg/logger"
	"github.com/web3sphere/backend/pkg/mailer"
	"github.com/web3sphere/backend/pkg/queue/kafka"
	"github.com/web3sphere/backend/pkg/queue/rabbitmq"
	"github.com/web3sphere/backend/pkg/storage"
	"github.com/web3sphere/backend/pkg/validator"
)

func main() {
	// 1. Load Configuration
	cfg := configs.Load()

	// 2. Initialize Logger
	log := logger.New(cfg.App.Env, cfg.App.Debug)
	defer log.Sync()
	log.Infof("Starting %s v%s in %s mode", cfg.App.Name, cfg.App.Version, cfg.App.Env)

	// 3. Initialize Database
	db, err := database.New(&cfg.Database, log)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close(db, log)

	// 4. Initialize Redis
	redisClient, err := cache.New(&cfg.Redis, log)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	defer redisClient.Close()

	// 5. Initialize RabbitMQ
	rmqClient, err := rabbitmq.New(&cfg.RabbitMQ, log)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer rmqClient.Close()

	// 6. Initialize Kafka (Infrastructure only)
	kafkaProducer, err := kafka.NewProducer(&cfg.Kafka, log)
	if err != nil {
		log.Warnf("Failed to initialize Kafka Producer: %v (continuing without Kafka)", err)
	} else {
		defer kafkaProducer.Close()
	}

	// 7. Initialize Mailer
	mailSvc, err := mailer.New(cfg, log)
	if err != nil {
		log.Fatalf("Failed to initialize Mailer: %v", err)
	}

	// 8. Initialize Storage
	storageSvc, err := storage.New(cfg, log)
	if err != nil {
		log.Fatalf("Failed to initialize Storage: %v", err)
	}
	_ = storageSvc // In a real app, this would be injected into services that need it

	// 9. Initialize JWT Manager
	jwtManager := jwt.NewManager(&cfg.JWT)

	// 10. Setup Gin Router
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	// Setup custom validators
	validator.Setup()

	// 11. Apply Global Middlewares
	router.Use(middleware.RequestID())
	router.Use(middleware.RequestLogger(log))
	router.Use(middleware.Recovery(log))
	router.Use(middleware.CORS(&cfg.CORS))
	router.Use(middleware.Compression())
	router.Use(middleware.IPDetection())
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.VersionHeader(cfg.App.Version))
	router.Use(middleware.Timeout(cfg.Server.ReadTimeout))

	// Base API Group
	api := router.Group("/api/v1")

	// Health Check
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "version": cfg.App.Version})
	})

	// 12. Initialize Modules
	auth.Setup(api, db, cfg, redisClient, jwtManager, mailSvc, log)
	users.Setup(api, db, jwtManager, redisClient, log)

	// 13. Setup HTTP Server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// 14. Graceful Shutdown
	go func() {
		log.Infof("Server listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Listen error: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Info("Server exiting")
}
