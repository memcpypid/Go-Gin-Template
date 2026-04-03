package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-gin-template/internal/config"
	"go-gin-template/internal/delivery/http/handler"
	"go-gin-template/internal/delivery/http/route"
	"go-gin-template/internal/middleware"
	"go-gin-template/internal/repository"
	"go-gin-template/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 1. Initialize Configuration (Viper & Config structs)
	v, err := config.NewViper(".")
	if err != nil {
		log.Fatalf("Failed to initialize viper: %v", err)
	}

	cfg, err := config.NewConfig(v)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Initialize Logger (Zap)
	logger, err := config.NewLogger(cfg.App.Env)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 3. Initialize Database (GORM)
	db, err := config.NewDatabase(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// 4. Initialize Validator & Translator
	validate, trans := config.NewValidator()

	// 5. Initialize Repositories
	userRepo := repository.NewUserRepository(db, logger)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db, logger)

	// 6. Initialize Services
	userService := service.NewUserService(userRepo, logger)
	authService := service.NewAuthService(userRepo, refreshTokenRepo, cfg, logger)

	// 7. Initialize Handlers
	userHandler := handler.NewUserHandler(userService, logger, trans, validate)
	authHandler := handler.NewAuthHandler(authService, logger, trans, validate)

	// 8. Initialize Middlewares
	mw := middleware.NewMiddleware(cfg, logger)

	// 9. New Router
	routeHandler := route.NewRouter(mw, authHandler, userHandler)
	engine := routeHandler.Setup()

	// 10. Start Server
	go func() {
		addr := fmt.Sprintf(":%d", cfg.App.Port)
		logger.Info(fmt.Sprintf("Starting server on %s", addr))
		if err := engine.Run(addr); err != nil {
			logger.Fatal("Failed to run server", zap.Error(err))
		}
	}()

	// 10. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}

	logger.Info("Server exiting")
}
