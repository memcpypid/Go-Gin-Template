package main

import (
	"log"

	"go-gin-template/internal/config"
	"go-gin-template/internal/entity"
	"go-gin-template/internal/utils"

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

	logger.Info("Running database seeder...")

	var count int64
	db.Model(&entity.User{}).Where("role = ?", "admin").Count(&count)

	if count == 0 {
		hashedPassword, err := utils.HashPassword(cfg.Admin.Password)
		if err != nil {
			logger.Fatal("Failed to hash default admin password", zap.Error(err))
		}

		admin := entity.User{
			Name:       "Admin",
			Email:      cfg.Admin.Email,
			Password:   hashedPassword,
			Role:       "admin",
			IsVerified: true,
		}
		
		if err := db.Create(&admin).Error; err != nil {
			logger.Fatal("Failed to seed admin user", zap.Error(err))
		}
		
		logger.Info("Successfully seeded default admin account.", zap.String("email", cfg.Admin.Email))
	} else {
		logger.Info("Admin account already exists. Seeding skipped.")
	}
}
