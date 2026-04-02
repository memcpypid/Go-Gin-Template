package service

import (
	"context"

	"go-gin-template/internal/dto"

	"github.com/google/uuid"
)

type UserService interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateProfileRequest) (*dto.UserResponse, error)
	GetUsers(ctx context.Context, limit, offset int, sort, sortBy string) ([]dto.UserResponse, int64, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
