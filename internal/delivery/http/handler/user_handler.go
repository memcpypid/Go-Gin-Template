package handler

import (
	"net/http"

	"go-gin-template/internal/dto"
	"go-gin-template/internal/service"
	"go-gin-template/internal/utils"
	"go-gin-template/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService service.UserService
	logger      *zap.Logger
}

func NewUserHandler(userService service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid user id"))
		return
	}

	userResponse, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("profile retrieved successfully", userResponse))
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid user id"))
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	userResponse, err := h.userService.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("profile updated successfully", userResponse))
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	pagination := utils.GeneratePaginationFromRequest(c)

	users, total, err := h.userService.GetUsers(c.Request.Context(), pagination.Limit, pagination.GetOffset(), pagination.Sort, pagination.SortBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	data := map[string]interface{}{
		"users": users,
		"total": total,
		"limit": pagination.Limit,
		"page":  pagination.Page,
	}

	c.JSON(http.StatusOK, response.Success("users retrieved successfully", data))
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid user id format"))
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("user deleted successfully", nil))
}

func (h *UserHandler) ActivateAccount(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		idStr = c.GetString("userID")
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid user id format"))
		return
	}

	if err := h.userService.ActivateAccount(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("account activated successfully", nil))
}

func (h *UserHandler) DeactivateAccount(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		idStr = c.GetString("userID")
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid user id format"))
		return
	}

	if err := h.userService.DeactivateAccount(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("account deactivated successfully", nil))
}
