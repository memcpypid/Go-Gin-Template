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
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid user id"))
		return
	}

	userResponse, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("profile retrieved successfully", userResponse))
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid user id"))
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ValidationError(err))
		return
	}

	userResponse, err := h.userService.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("profile updated successfully", userResponse))
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	pagination := utils.GeneratePaginationFromRequest(c)

	users, total, err := h.userService.GetUsers(c.Request.Context(), pagination.Limit, pagination.GetOffset(), pagination.Search, pagination.Sort, pagination.SortBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithPagination(
		"users retrieved successfully",
		users,
		total,
		pagination.Limit,
		pagination.Page,
	))
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid user id format"))
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, err.Error()))
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
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid user id format"))
		return
	}

	if err := h.userService.ActivateAccount(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, err.Error()))
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
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid user id format"))
		return
	}

	if err := h.userService.DeactivateAccount(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("account deactivated successfully", nil))
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid user id format"))
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ValidationError(err))
		return
	}

	userResponse, err := h.userService.UpdateUser(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("user updated successfully", userResponse))
}
