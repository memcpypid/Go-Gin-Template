package middleware

import (
	"go-gin-template/internal/config"
	"go.uber.org/zap"
)

type Middleware struct {
	cfg    *config.Config
	logger *zap.Logger
}

func NewMiddleware(cfg *config.Config, logger *zap.Logger) *Middleware {
	return &Middleware{
		cfg:    cfg,
		logger: logger,
	}
}
