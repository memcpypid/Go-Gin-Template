package repository

import (
	"context"
	"errors"

	"go-gin-template/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"go.uber.org/zap"
)

type userRepositoryImpl struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewUserRepository(db *gorm.DB, logger *zap.Logger) UserRepository {
	return &userRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
	r.logger.Info("Repository: Create user called", zap.String("email", user.Email))
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil when not found to easily check
		}
		r.logger.Error("Repository: GetByID failed", zap.Error(err), zap.String("id", id.String()))
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Error("Repository: GetByEmail failed", zap.Error(err), zap.String("email", email))
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search, sort, sortBy string) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	db := r.db.WithContext(ctx).Model(&entity.User{})

	if search != "" {
		searchTerm := "%" + search + "%"
		db = db.Where("name ILIKE ? OR email ILIKE ?", searchTerm, searchTerm)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if sortBy == "" {
		sortBy = "created_at"
	}
	if sort != "asc" && sort != "desc" {
		sort = "desc"
	}
	orderClause := sortBy + " " + sort

	if err := db.Limit(limit).Offset(offset).Order(orderClause).Find(&users).Error; err != nil {
		r.logger.Error("Repository: FindAll failed", zap.Error(err))
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepositoryImpl) Update(ctx context.Context, user *entity.User) error {
	r.logger.Info("Repository: Update user called", zap.String("id", user.ID.String()))
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	r.logger.Info("Repository: Delete user called", zap.String("id", id.String()))
	return r.db.WithContext(ctx).Delete(&entity.User{}, "id = ?", id).Error
}
