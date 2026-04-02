package service

import (
	"context"
	"errors"

	"go-gin-template/internal/config"
	"go-gin-template/internal/dto"
	"go-gin-template/internal/entity"
	"go-gin-template/internal/repository"
	"go-gin-template/internal/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type authServiceImpl struct {
	userRepo repository.UserRepository
	cfg      *config.Config
	logger   *zap.Logger
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config, logger *zap.Logger) AuthService {
	return &authServiceImpl{
		userRepo: userRepo,
		cfg:      cfg,
		logger:   logger,
	}
}

func (s *authServiceImpl) Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error) {
	s.logger.Info("Service: Register called", zap.String("email", req.Email))

	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.logger.Error("Service: Failed to hash password", zap.Error(err))
		return nil, err
	}

	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user", 
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	res := dto.ToUserResponse(user)
	return &res, nil
}

func (s *authServiceImpl) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	s.logger.Info("Service: Login called", zap.String("email", req.Email))

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	access, refresh, err := utils.GenerateTokens(
		user.ID,
		user.Role,
		s.cfg.JWT.Secret,
		s.cfg.JWT.ExpirationHours,
		s.cfg.JWT.RefreshExpirationHours,
	)
	if err != nil {
		s.logger.Error("Service: Failed to generate tokens", zap.Error(err))
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		User:         dto.ToUserResponse(user),
	}, nil
}

func (s *authServiceImpl) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.TokenResponse, error) {
	claims, err := utils.ValidateToken(req.RefreshToken, s.cfg.JWT.Secret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	userIDStr := claims.Subject
	if userIDStr == "" {
		userIDStr = claims.UserID.String()
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user id in token")
	}

	s.logger.Info("Service: RefreshToken called", zap.String("user_id", userID.String()))

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	access, refresh, err := utils.GenerateTokens(
		user.ID,
		user.Role,
		s.cfg.JWT.Secret,
		s.cfg.JWT.ExpirationHours,
		s.cfg.JWT.RefreshExpirationHours,
	)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
