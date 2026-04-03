package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Admin    AdminConfig
}

type AppConfig struct {
	Env  string `mapstructure:"APP_ENV"`
	Port int    `mapstructure:"APP_PORT"`
	Name string `mapstructure:"APP_NAME"`
}

type DatabaseConfig struct {
	Driver  string `mapstructure:"DB_DRIVER"`
	Host    string `mapstructure:"DB_HOST"`
	Port    int    `mapstructure:"DB_PORT"`
	User    string `mapstructure:"DB_USER"`
	Pass    string `mapstructure:"DB_PASS"`
	Name    string `mapstructure:"DB_NAME"`
	SSLMode string `mapstructure:"DB_SSLMODE"`
}

type JWTConfig struct {
	Secret            string `mapstructure:"JWT_SECRET"`
	Expiration        string `mapstructure:"JWT_EXPIRATION"`
	RefreshExpiration string `mapstructure:"JWT_REFRESH_EXPIRATION"`
}

type AdminConfig struct {
	Email    string `mapstructure:"ADMIN_EMAIL"`
	Password string `mapstructure:"ADMIN_PASSWORD"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path + "/.env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading env file, falling back to environment variables: %v", err)
	}

	var cfg Config

	viper.BindEnv("APP_ENV")
	viper.BindEnv("APP_PORT")
	viper.BindEnv("APP_NAME")
	viper.BindEnv("DB_DRIVER")
	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_PORT")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASS")
	viper.BindEnv("DB_NAME")
	viper.BindEnv("DB_SSLMODE")
	viper.BindEnv("JWT_SECRET")
	viper.BindEnv("JWT_EXPIRATION")
	viper.BindEnv("JWT_REFRESH_EXPIRATION")
	viper.BindEnv("ADMIN_EMAIL")
	viper.BindEnv("ADMIN_PASSWORD")

	if err := viper.Unmarshal(&cfg.App); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&cfg.Database); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&cfg.JWT); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&cfg.Admin); err != nil {
		return nil, err
	}

	if cfg.App.Port == 0 {
		cfg.App.Port = 8080 // Default
	}
	if cfg.JWT.Secret == "" {
		cfg.JWT.Secret = "default_secret_key"
	}
	if cfg.JWT.Expiration == "" {
		cfg.JWT.Expiration = "24h"
	}
	if cfg.JWT.RefreshExpiration == "" {
		cfg.JWT.RefreshExpiration = "168h"
	}

	return &cfg, nil
}
