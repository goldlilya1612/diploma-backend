package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name  string    `gorm:"type:varchar(255);not null"`
	Email string    `gorm:"uniqueIndex;not null"`
	Role  string    `gorm:"type:varchar(255);notnull"`
	//Provider  string    `gorm:"not null"`
	Verified bool `gorm:"not null"`

	CreateAt time.Time
	UpdateAt time.Time
}

type SignUp struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
	Photo           string `json:"photo" binding:"required"`
}

type SignInInput struct {
	Email    string `json:"email"  binding:"required,email"`
	Password string `json:"password"  binding:"required"`
}

type UserResponse struct {
	ID    uuid.UUID `json:"id,omitempty"`
	Name  string    `json:"name,omitempty"`
	Email string    `json:"email,omitempty"`
	Role  string    `json:"role,omitempty"`
	//Provider  string    `json:"provider"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
