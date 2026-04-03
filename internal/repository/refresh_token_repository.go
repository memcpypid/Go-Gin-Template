package repository

import (
	"context"
	"go-gin-template/internal/entity"

	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	BaseRepository[entity.RefreshToken]
	Create(ctx context.Context, token *entity.RefreshToken) error
	GetByToken(ctx context.Context, token string) (*entity.RefreshToken, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}
