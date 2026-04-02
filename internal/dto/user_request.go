package dto

type UpdateProfileRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password,omitempty" binding:"omitempty,min=6"`
}
