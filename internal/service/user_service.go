package service

import (
	"context"

	"go-gin-template/internal/dto"

	"github.com/google/uuid"
)

type UserService interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateProfileRequest) (*dto.UserResponse, error)
	GetUsers(ctx context.Context, limit, offset int, search, sort, sortBy string) ([]dto.UserResponse, int64, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	ActivateAccount(ctx context.Context, id uuid.UUID) error
	DeactivateAccount(ctx context.Context, id uuid.UUID) error
}
