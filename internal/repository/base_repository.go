package repository

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BaseRepository[T any] interface {
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*T, error)
	Count(ctx context.Context) (int64, error)
}

type baseRepositoryImpl[T any] struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewBaseRepository[T any](db *gorm.DB, logger *zap.Logger) BaseRepository[T] {
	return &baseRepositoryImpl[T]{
		db:     db,
		logger: logger,
	}
}

func (r *baseRepositoryImpl[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *baseRepositoryImpl[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *baseRepositoryImpl[T]) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(new(T), "id = ?", id).Error
}

func (r *baseRepositoryImpl[T]) GetByID(ctx context.Context, id uuid.UUID) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (r *baseRepositoryImpl[T]) Count(ctx context.Context) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(new(T)).Count(&total).Error
	return total, err
}
