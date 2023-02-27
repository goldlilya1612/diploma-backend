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
	Image       string    `gorm:"type:varchar(255)"`
	Category    string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text;not null"`
	//Chapters    Chapters

	CreateAt time.Time
	UpdateAt time.Time

	Lecturer Lecturer `gorm:"foreignKey:CreatorID;references:UserID"`
}

type CreateCourse struct {
	Name        string    `json:"name" binding:"required"`
	CreatorID   uuid.UUID `json:"creatorID" binding:"required"`
	CreatorName string    `json:"creatorName" binding:"required"`
	Image       string    `json:"image" binding:"required"`
	Category    string    `json:"category" binding:"required"`
	Description string    `json:"description" binding:"required"`
}

type CourseResponse struct {
	ID          int       `json:"id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	CreatorID   uuid.UUID `json:"creatorID" binding:"required"`
	CreatorName string    `json:"creatorName" binding:"required"`
	Image       string    `json:"image" binding:"required"`
	Category    string    `json:"category" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Route       string    `json:"route" binding:"required"`

	CreateAt time.Time `json:"createAt"`
	UpdateAt time.Time `json:"updateAt"`
}

type Chapters struct {
}
