package models

import (
	"github.com/google/uuid"
	"time"
)

type Course struct {
	ID          int       `gorm:"primaryKey;uniqueIndex"`
	Name        string    `gorm:"type:varchar(255);not null"`
	CreatorID   uuid.UUID `gorm:"type:uuid;not null"`
	CreatorName string    `gorm:"type:varchar(255);not null"`
	Image       int       `gorm:"type:integer"`
	Category    string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Images   Images   `gorm:"foreignKey:Image;constraint:OnDelete:CASCADE;"`
	Lecturer Lecturer `gorm:"foreignKey:CreatorID;references:UserID"`
}

type CreateCourse struct {
	Name string `form:"name" binding:"required"`
	//Image       multipart.File `form:"image,omitempty" validate:"required"`
	Category    string `form:"category" binding:"required"`
	Description string `form:"description" binding:"required"`
}

type UpdateCourse struct {
	ID   int    `form:"id" binding:"required"`
	Name string `form:"name" binding:"required"`
	//Image       multipart.File `form:"image,omitempty" validate:"required"`
	Category    string `form:"category" binding:"required"`
	Description string `form:"description" binding:"required"`
}

type CourseResponse struct {
	ID          int       `json:"id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	CreatorID   uuid.UUID `json:"creatorID" binding:"required"`
	CreatorName string    `json:"creatorName" binding:"required"`
	Image       Images    `json:"image" binding:"required"`
	Category    string    `json:"category" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Route       string    `json:"route" binding:"required"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
