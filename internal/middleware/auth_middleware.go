package middleware

import (
	"net/http"
	"strings"

	"go-gin-template/internal/config"
	"go-gin-template/internal/utils"
	"go-gin-template/pkg/response"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, response.Error("authorization header is required"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, response.Error("authorization header format must be Bearer {token}"))
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString, cfg.JWT.Secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.Error("invalid or expired token"))
			c.Abort()
			return
		}

		// Set context
		c.Set("userID", claims.UserID.String())
		c.Set("role", claims.Role)
		c.Next()
	}
}
