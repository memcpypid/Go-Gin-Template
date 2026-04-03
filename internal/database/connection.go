package database

import (
	"fmt"

	"go-gin-template/internal/config"
	"go-gin-template/internal/entity"
	"go-gin-template/internal/utils"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDatabase(cfg *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector
	var dsn string

	if cfg.Database.Driver == "postgres" {
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Pass, cfg.Database.Name, cfg.Database.SSLMode)
		dialector = postgres.Open(dsn)
	} else if cfg.Database.Driver == "mysql" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Database.User, cfg.Database.Pass, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
		dialector = mysql.Open(dsn)
	} else {
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Database.Driver)
	}

	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	// Connection pooling tuning could go here
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// Automigrate
	err = gormDB.AutoMigrate(&entity.User{}, &entity.RefreshToken{})
	if err != nil {
		return nil, err
	}

	// Seed Default Admin
	seedAdmin(gormDB, cfg)

	return gormDB, nil
}

func seedAdmin(db *gorm.DB, cfg *config.Config) {
	var count int64
	db.Model(&entity.User{}).Where("role = ?", "admin").Count(&count)

	if count == 0 {
		hashedPassword, err := utils.HashPassword(cfg.Admin.Password)
		if err != nil {
			panic("failed to hash default admin password")
		}

		admin := entity.User{
			Name:     "Admin",
			Email:    cfg.Admin.Email,
			Password: hashedPassword,
			Role:     "admin",
		}
		db.Create(&admin)
	}
}
