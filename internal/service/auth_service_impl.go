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
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	cfg              *config.Config
	logger           *zap.Logger
}

func NewAuthService(userRepo repository.UserRepository, refreshTokenRepo repository.RefreshTokenRepository, cfg *config.Config, logger *zap.Logger) AuthService {
	return &authServiceImpl{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		cfg:              cfg,
		logger:           logger,
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

	if !user.IsVerified {
		return nil, errors.New("account is not verified")
	}

	accessDuration, err := utils.ParseDuration(s.cfg.JWT.Expiration)
	if err != nil {
		s.logger.Error("Service: Invalid access token duration", zap.Error(err))
		return nil, err
	}
	refreshDuration, err := utils.ParseDuration(s.cfg.JWT.RefreshExpiration)
	if err != nil {
		s.logger.Error("Service: Invalid refresh token duration", zap.Error(err))
		return nil, err
	}

	access, refresh, refreshExpAt, err := utils.GenerateTokens(
		user.ID,
		user.Role,
		s.cfg.JWT.Secret,
		accessDuration,
		refreshDuration,
	)
	if err != nil {
		s.logger.Error("Service: Failed to generate tokens", zap.Error(err))
		return nil, err
	}

	// Save refresh token to DB
	refreshTokenEntity := &entity.RefreshToken{
		Token:     refresh,
		UserID:    user.ID,
		ExpiresAt: refreshExpAt,
	}
	if err := s.refreshTokenRepo.Create(ctx, refreshTokenEntity); err != nil {
		s.logger.Error("Service: Failed to save refresh token", zap.Error(err))
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		User:         dto.ToUserResponse(user),
	}, nil
}

func (s *authServiceImpl) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.TokenResponse, error) {
	// Validate token in DB
	storedToken, err := s.refreshTokenRepo.GetByToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	if storedToken == nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// Validate JWT
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

	accessDuration, err := utils.ParseDuration(s.cfg.JWT.Expiration)
	if err != nil {
		return nil, err
	}
	refreshDuration, err := utils.ParseDuration(s.cfg.JWT.RefreshExpiration)
	if err != nil {
		return nil, err
	}

	access, refresh, refreshExpAt, err := utils.GenerateTokens(
		user.ID,
		user.Role,
		s.cfg.JWT.Secret,
		accessDuration,
		refreshDuration,
	)
	if err != nil {
		return nil, err
	}

	// Revoke old token
	if err := s.refreshTokenRepo.Revoke(ctx, storedToken.ID); err != nil {
		s.logger.Error("Service: Failed to revoke old refresh token", zap.Error(err))
		// Continue even if revoking fails, but maybe log it
	}

	// Save new refresh token
	newRefreshTokenEntity := &entity.RefreshToken{
		Token:     refresh,
		UserID:    user.ID,
		ExpiresAt: refreshExpAt,
	}
	if err := s.refreshTokenRepo.Create(ctx, newRefreshTokenEntity); err != nil {
		s.logger.Error("Service: Failed to save new refresh token", zap.Error(err))
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (s *authServiceImpl) Logout(ctx context.Context, req dto.LogoutRequest) error {
	s.logger.Info("Service: Logout called")

	// Validate token in DB
	storedToken, err := s.refreshTokenRepo.GetByToken(ctx, req.RefreshToken)
	if err != nil {
		return err
	}
	if storedToken == nil {
		return errors.New("invalid or already revoked refresh token")
	}

	// Revoke token
	if err := s.refreshTokenRepo.Revoke(ctx, storedToken.ID); err != nil {
		s.logger.Error("Service: Failed to revoke refresh token", zap.Error(err))
		return err
	}

	return nil
}
