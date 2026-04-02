package test

import (
	"context"
	"errors"
	"testing"

	"go-gin-template/internal/entity"
	"go-gin-template/internal/service"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// A tiny stub repository for testing
type stubUserRepo struct {
	MockGetByID func(ctx context.Context, id uuid.UUID) (*entity.User, error)
}

func (s *stubUserRepo) Create(ctx context.Context, user *entity.User) error { return nil }
func (s *stubUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	return s.MockGetByID(ctx, id)
}
func (s *stubUserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	return nil, nil
}
func (s *stubUserRepo) FindAll(ctx context.Context, limit, offset int, sort, sortBy string) ([]entity.User, int64, error) {
	return nil, 0, nil
}
func (s *stubUserRepo) Update(ctx context.Context, user *entity.User) error { return nil }
func (s *stubUserRepo) Delete(ctx context.Context, id uuid.UUID) error { return nil }

func TestGetProfile_Success(t *testing.T) {
	id := uuid.New()
	repo := &stubUserRepo{
		MockGetByID: func(ctx context.Context, reqId uuid.UUID) (*entity.User, error) {
			if id == reqId {
				return &entity.User{ID: id, Name: "Test User", Email: "test@example.com"}, nil
			}
			return nil, errors.New("wrong id")
		},
	}
	logger := zap.NewNop()
	svc := service.NewUserService(repo, logger)

	profile, err := svc.GetProfile(context.Background(), id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if profile.Name != "Test User" {
		t.Fatalf("expected Test User, got %v", profile.Name)
	}
}

func TestGetProfile_NotFound(t *testing.T) {
	repo := &stubUserRepo{
		MockGetByID: func(ctx context.Context, reqId uuid.UUID) (*entity.User, error) {
			return nil, nil // user nil is not found
		},
	}
	logger := zap.NewNop()
	svc := service.NewUserService(repo, logger)

	_, err := svc.GetProfile(context.Background(), uuid.New())
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if err.Error() != "user not found" {
		t.Fatalf("expected user not found error, got %v", err)
	}
}
