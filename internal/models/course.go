package models

import (
	"github.com/google/uuid"
	"time"
)

type Course struct {
	ID          int       `gorm:"primaryKey;uniqueIndex"`
	Name        string    `gorm:"type:varchar(255);not null"`
	CreatorID   uuid.UUID `gorm:"type:uuid;not null"`
	ImageID     int       `gorm:"type:integer"`
	Category    string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Image    Image    `gorm:"foreignKey:ImageID;constraint:OnDelete:CASCADE;"`
	Lecturer Lecturer `gorm:"foreignKey:CreatorID;references:UserID"`

	Chapters []Chapter
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
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CreatorName string `json:"creatorName"`
	Image       Image  `json:"image"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Route       string `json:"route"`

	Chapters []Chapter `json:"chapters"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
