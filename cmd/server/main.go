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
	"github.com/web3sphere/backend/internal/auth"
	"github.com/web3sphere/backend/internal/bootstrap"
	"github.com/web3sphere/backend/internal/middleware"
	"github.com/web3sphere/backend/internal/users"
	"github.com/web3sphere/backend/pkg/validator"
)

func main() {
	// 1. Initialize Global Container (DB, Redis, Logger, Config, etc.)
	container, cleanup, err := bootstrap.InitContainer()
	if err != nil {
		fmt.Printf("Fatal error initializing container: %v\n", err)
		os.Exit(1)
	}
	defer cleanup() // Automatically closes all connections on exit

	// 2. Setup Gin Router
	if container.Config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	// 3. Setup Custom Validators
	validator.Setup()

	// 4. Apply Global Middlewares cleanly
	middleware.ApplyGlobal(router, container.Config, container.Logger)

	// Base API Group
	api := router.Group("/api/v1")

	// Health Check
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "version": container.Config.App.Version})
	})

	// 5. Initialize Modules efficiently using the container
	auth.Setup(api, container.DB, container.Config, container.Redis, container.JWTManager, container.Mailer, container.Logger)
	users.Setup(api, container.DB, container.JWTManager, container.Redis, container.Logger)
	
	// Example of future modules:
	// projects.Setup(api, container.DB, container.Redis, container.Logger)
	// escrow.Setup(api, container.DB, container.Redis, container.Logger)

	// 6. Setup HTTP Server
	addr := fmt.Sprintf("%s:%s", container.Config.Server.Host, container.Config.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  container.Config.Server.ReadTimeout,
		WriteTimeout: container.Config.Server.WriteTimeout,
		IdleTimeout:  container.Config.Server.IdleTimeout,
	}

	// 7. Graceful Shutdown listener
	go func() {
		container.Logger.Infof("Server listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			container.Logger.Fatalf("Listen error: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	container.Logger.Info("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), container.Config.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		container.Logger.Fatalf("Server forced to shutdown: %v", err)
	}

	container.Logger.Info("Server exiting")
}
