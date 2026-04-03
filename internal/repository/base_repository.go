package repository

import (
	"context"

	"go-gin-template/internal/utils"

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
	BuildPaginationQuery(db *gorm.DB, pagination *utils.Pagination, searchFields []string) *gorm.DB
	Paginate(pagination *utils.Pagination) func(db *gorm.DB) *gorm.DB
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

// Paginate generates the DB Scope function for limits and offsets
func (r *baseRepositoryImpl[T]) Paginate(pagination *utils.Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		orderStr := pagination.SortBy + " " + pagination.Sort
		return db.Offset(pagination.GetOffset()).Limit(pagination.Limit).Order(orderStr)
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

func (r *baseRepositoryImpl[T]) BuildPaginationQuery(db *gorm.DB, pagination *utils.Pagination, searchFields []string) *gorm.DB {
	query := db.Model(new(T))

	if pagination.Search != "" && len(searchFields) > 0 {
		searchQuery := ""
		likeOperator := "LIKE"
		if db.Dialector.Name() == "postgres" {
			likeOperator = "ILIKE"
		}

		var values []interface{}
		for i, field := range searchFields {
			if i > 0 {
				searchQuery += " OR "
			}
			searchQuery += field + " " + likeOperator + " ?"
			values = append(values, "%"+pagination.Search+"%")
		}
		query = query.Where(searchQuery, values...)
	}

	return query
}
