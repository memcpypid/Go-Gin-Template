package repository

import (
	"context"
	"errors"

	"go-gin-template/internal/entity"
	"go-gin-template/internal/utils"

	"go.uber.org/zap"
	"gorm.io/gorm"
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
	r.logger.Info("Repository: Getting user by email", zap.String("email", email))
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

func (r *userRepositoryImpl) FindAll(ctx context.Context, p *utils.Pagination) ([]entity.User, int64, error) {
	r.logger.Info("Repository: Finding all users")
	var users []entity.User
	var total int64

	// Build query with search logic using helper
	query := r.BuildPaginationQuery(r.db.WithContext(ctx), p, []string{"name", "email"})

	// Get total count before pagination
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Repository: FindAll count failed", zap.Error(err))
		return nil, 0, err
	}

	// Execute paginated query
	if err := query.Scopes(r.Paginate(p)).Find(&users).Error; err != nil {
		r.logger.Error("Repository: FindAll query failed", zap.Error(err))
		return nil, 0, err
	}

	return users, total, nil
}
