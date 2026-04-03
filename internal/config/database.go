package config

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"moul.io/zapgorm2"
)

func NewDatabase(cfg *Config, logger *zap.Logger) (*gorm.DB, error) {
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

	zapGormLogger := zapgorm2.New(logger)
	zapGormLogger.SetAsDefault()

	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: zapGormLogger.LogMode(gormLogger.Info),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	// Connection pooling tuning
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return gormDB, nil
}
