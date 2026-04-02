package http

import (
	"go-gin-template/internal/config"
	"go-gin-template/internal/delivery/http/handler"
	"go-gin-template/internal/middleware"
	"go-gin-template/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupRouter(
	cfg *config.Config,
	logger *zap.Logger,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
) *gin.Engine {

	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggingMiddleware(logger))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, response.Success("server is healthy", nil))
	})

	api := router.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		protected.GET("/me", userHandler.GetProfile)
		protected.PUT("/me", userHandler.UpdateProfile)

		admin := protected.Group("/users")
		admin.Use(middleware.RoleMiddleware("admin"))
		{
			admin.GET("", userHandler.GetUsers)
			admin.DELETE("/:id", userHandler.DeleteUser)
			admin.PATCH("/:id/activate", userHandler.ActivateAccount)
			admin.PATCH("/:id/deactivate", userHandler.DeactivateAccount)
		}
	}

	return router
}
