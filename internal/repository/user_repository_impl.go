package repository

import (
	"context"
	"errors"

	"go-gin-template/internal/entity"

	"gorm.io/gorm"
	"go.uber.org/zap"
)

type userRepositoryImpl struct {
	BaseRepository[entity.User]
	db     *gorm.DB
	logger *zap.Logger
}

func NewUserRepository(db *gorm.DB, logger *zap.Logger) UserRepository {
	return &userRepositoryImpl{
		BaseRepository: NewBaseRepository[entity.User](db, logger),
		db:             db,
		logger:         logger,
	}
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
