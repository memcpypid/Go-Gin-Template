package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggingMiddleware(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		logFields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("latency", latency),
		}

		if len(c.Errors) > 0 {
			logFields = append(logFields, zap.String("errors", c.Errors.String()))
			log.Error("Request failed", logFields...)
		} else if status >= 400 && status < 500 {
			log.Warn("Client error", logFields...)
		} else if status >= 500 {
			log.Error("Server error", logFields...)
		} else {
			log.Info("Request handled", logFields...)
		}
	}
}
