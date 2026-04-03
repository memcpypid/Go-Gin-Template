package http

import (
	"go-gin-template/internal/delivery/http/handler"
	"go-gin-template/internal/middleware"
	"go-gin-template/pkg/response"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	mw *middleware.Middleware,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
) *gin.Engine {

	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(mw.CORSMiddleware())
	router.Use(mw.LoggingMiddleware())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, response.Success("server is healthy", nil))
	})

	api := router.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.Logout)
	}

	protected := api.Group("/")
	protected.Use(mw.AuthMiddleware())
	{
		protected.GET("/users/me", userHandler.GetProfile)
		protected.PUT("/users/me", userHandler.UpdateProfile)

		admin := protected.Group("/users")
		admin.Use(mw.RoleMiddleware("admin"))
		{
			admin.GET("", userHandler.GetUsers)
			admin.PUT("/:id", userHandler.UpdateUser)
			admin.DELETE("/:id", userHandler.DeleteUser)
			admin.PATCH("/:id/activate", userHandler.ActivateAccount)
			admin.PATCH("/:id/deactivate", userHandler.DeactivateAccount)
		}
	}

	return router
}
