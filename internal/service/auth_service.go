package service

import (
	"context"

	"go-gin-template/internal/dto"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.TokenResponse, error)
	Logout(ctx context.Context, req dto.LogoutRequest) error
}
