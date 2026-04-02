package handler

import (
	"net/http"

	"go-gin-template/internal/dto"
	"go-gin-template/internal/service"
	"go-gin-template/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService service.AuthService
	logger      *zap.Logger
}

func NewAuthHandler(authService service.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	userResponse, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.Success("user registered successfully", userResponse))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	loginResponse, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("login successful", loginResponse))
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	tokenResponse, err := h.authService.RefreshToken(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("token refreshed successfully", tokenResponse))
}
