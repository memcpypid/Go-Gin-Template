package route

import (
	"go-gin-template/internal/delivery/http/handler"
	"go-gin-template/internal/middleware"
	"go-gin-template/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	mw          *middleware.Middleware
	authHandler *handler.AuthHandler
	userHandler *handler.UserHandler
}

func NewRouter(
	mw *middleware.Middleware,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
) *Router {
	return &Router{
		mw:          mw,
		authHandler: authHandler,
		userHandler: userHandler,
	}
}

func (r *Router) Setup() *gin.Engine {
	engine := gin.New()

	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())
	engine.Use(r.mw.CORSMiddleware())
	engine.Use(r.mw.LoggingMiddleware())

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, response.Success(http.StatusOK, "server is healthy", nil))
	})

	api := engine.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/login", r.authHandler.Login)
		auth.POST("/refresh", r.authHandler.RefreshToken)
		auth.POST("/logout", r.authHandler.Logout)
	}

	protected := api.Group("/")
	protected.Use(r.mw.AuthMiddleware())
	{
		protected.GET("/users/me", r.userHandler.GetProfile)
		protected.PUT("/users/me", r.userHandler.UpdateProfile)

		admin := protected.Group("/users")
		admin.Use(r.mw.RoleMiddleware("admin"))
		{
			admin.GET("", r.userHandler.GetUsers)
			admin.GET("/count", r.userHandler.GetCountUser)
			admin.PUT("/:id", r.userHandler.UpdateUser)
			admin.DELETE("/:id", r.userHandler.DeleteUser)
			admin.PATCH("/:id/activate", r.userHandler.ActivateAccount)
			admin.PATCH("/:id/deactivate", r.userHandler.DeactivateAccount)
		}
	}

	return engine
}
