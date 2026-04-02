package dto

type UpdateProfileRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password,omitempty" binding:"omitempty,min=6"`
}

type UpdateUserRequest struct {
	Name     string `json:"name,omitempty" binding:"omitempty"`
	Email    string `json:"email,omitempty" binding:"omitempty,email"`
	Password string `json:"password,omitempty" binding:"omitempty,min=6"`
	Role     string `json:"role,omitempty" binding:"omitempty,oneof=admin user"`
}
