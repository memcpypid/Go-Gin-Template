package middleware

import (
	"net/http"

	"go-gin-template/pkg/response"
	"github.com/gin-gonic/gin"
)

func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "role profile not found"))
			c.Abort()
			return
		}

		userRole := role.(string)
		found := false
		for _, r := range requiredRoles {
			if r == userRole {
				found = true
				break
			}
		}

		if !found {
			c.JSON(http.StatusForbidden, response.Error(http.StatusForbidden, "you don't have permission to access this resource"))
			c.Abort()
			return
		}

		c.Next()
	}
}
