package service

import (
	"context"
	"errors"

	"go-gin-template/internal/dto"
	"go-gin-template/internal/repository"
	"go-gin-template/internal/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type userServiceImpl struct {
	userRepo repository.UserRepository
	logger   *zap.Logger
}

func NewUserService(userRepo repository.UserRepository, logger *zap.Logger) UserService {
	return &userServiceImpl{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *userServiceImpl) GetProfile(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	s.logger.Info("Service: GetProfile called", zap.String("user_id", userID.String()))

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	res := dto.ToUserResponse(user)
	return &res, nil
}

func (s *userServiceImpl) UpdateProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	s.logger.Info("Service: UpdateProfile called", zap.String("user_id", userID.String()))

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	user.Name = req.Name
	if req.Password != "" {
		hash, err := utils.HashPassword(req.Password)
		if err != nil {
			s.logger.Error("Service: Failed to hash password", zap.Error(err))
			return nil, err
		}
		user.Password = hash
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	res := dto.ToUserResponse(user)
	return &res, nil
}

func (s *userServiceImpl) GetUsers(ctx context.Context, limit, offset int, sort, sortBy string) ([]dto.UserResponse, int64, error) {
	users, total, err := s.userRepo.FindAll(ctx, limit, offset, sort, sortBy)
	if err != nil {
		return nil, 0, err
	}

	return dto.ToUserResponseList(users), total, nil
}

func (s *userServiceImpl) DeleteUser(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("Service: DeleteUser called", zap.String("user_id", id.String()))

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	
	if user.Role == "admin" {
		return errors.New("cannot delete admin user")
	}

	return s.userRepo.Delete(ctx, id)
}

func (s *userServiceImpl) ActivateAccount(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("Service: ActivateAccount called", zap.String("user_id", id.String()))

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	user.IsVerified = true
	return s.userRepo.Update(ctx, user)
}

func (s *userServiceImpl) DeactivateAccount(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("Service: DeactivateAccount called", zap.String("user_id", id.String()))

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	user.IsVerified = false
	return s.userRepo.Update(ctx, user)
}
