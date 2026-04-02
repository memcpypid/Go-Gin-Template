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
			c.JSON(http.StatusUnauthorized, response.Error("unauthorized: role not found in context"))
			c.Abort()
			return
		}

		userRole := role.(string)
		isAllowed := false

		for _, reqRole := range requiredRoles {
			if userRole == reqRole {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.JSON(http.StatusForbidden, response.Error("forbidden: insufficient privileges"))
			c.Abort()
			return
		}

		c.Next()
	}
}
