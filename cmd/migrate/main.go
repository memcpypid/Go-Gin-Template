package main

import (
	"log"

	"go-gin-template/internal/config"
	"go-gin-template/internal/entity"

	"go.uber.org/zap"
)

func main() {
	v, err := config.NewViper(".")
	if err != nil {
		log.Fatalf("Failed to initialize viper: %v", err)
	}

	cfg, err := config.NewConfig(v)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger, err := config.NewLogger(cfg.App.Env)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	db, err := config.NewDatabase(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	logger.Info("Running database migrations...")

	err = db.AutoMigrate(
		&entity.User{},
		&entity.RefreshToken{},
	)

	if err != nil {
		logger.Fatal("Database migration failed", zap.Error(err))
	}

	logger.Info("Database migration completed successfully.")
}
