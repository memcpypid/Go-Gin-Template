package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;"`
	Name      string         `gorm:"type:varchar(100);not null"`
	Email     string         `gorm:"type:varchar(100);unique;not null"`
	Password   string         `gorm:"type:varchar(255);not null"`
	Role       string         `gorm:"type:varchar(20);not null;default:'user'"` // admin, user
	IsVerified bool           `gorm:"type:boolean;default:false"`
	CreatedAt  time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate hook to generate UUID
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}
