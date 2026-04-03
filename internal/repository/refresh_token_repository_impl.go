package repository

import (
	"context"
	"time"

	"go-gin-template/internal/entity"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type refreshTokenRepositoryImpl struct {
	BaseRepository[entity.RefreshToken]
	db     *gorm.DB
	logger *zap.Logger
}

func NewRefreshTokenRepository(db *gorm.DB, logger *zap.Logger) RefreshTokenRepository {
	return &refreshTokenRepositoryImpl{
		BaseRepository: NewBaseRepository[entity.RefreshToken](db, logger),
		db:             db,
		logger:         logger,
	}
}

func (r *refreshTokenRepositoryImpl) GetByToken(ctx context.Context, token string) (*entity.RefreshToken, error) {
	r.logger.Info("Repository: GetByToken RefreshToken called")
	var refreshToken entity.RefreshToken
	err := r.db.WithContext(ctx).Where("token = ? AND revoked_at IS NULL", token).First(&refreshToken).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &refreshToken, nil
}

func (r *refreshTokenRepositoryImpl) Revoke(ctx context.Context, id uuid.UUID) error {
	r.logger.Info("Repository: Revoke RefreshToken called", zap.String("id", id.String()))
	now := time.Now()
	return r.db.WithContext(ctx).Model(&entity.RefreshToken{}).Where("id = ?", id).Update("revoked_at", &now).Error
}

func (r *refreshTokenRepositoryImpl) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	r.logger.Info("Repository: DeleteByUserID RefreshToken called", zap.String("user_id", userID.String()))
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&entity.RefreshToken{}).Error
}
