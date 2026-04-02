package dto

import (
	"time"

	"github.com/google/uuid"
	"go-gin-template/internal/entity"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToUserResponse(user *entity.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUserResponseList(users []entity.User) []UserResponse {
	var responses []UserResponse
	if len(users) == 0 {
		return []UserResponse{}
	}
	for _, user := range users {
		responses = append(responses, ToUserResponse(&user))
	}
	return responses
}
