package models

import "time"

type Chapter struct {
	ID       int    `gorm:"primaryKey;uniqueIndex"`
	Name     string `gorm:"type:varchar(255);not null"`
	CourseID int    `gorm:"type:integer;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Course Course `gorm:"foreignKey:CourseID"`
}

type CreateChapter struct {
	Name     string `json:"name" binding:"required"`
	CourseID int    `json:"courseID" binding:"required"`
}

type UpdateChapter struct {
	ID   int    `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type ChapterResponse struct {
	ID       int    `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	CourseID int    `json:"courseID" binding:"required"`
	Route    string `json:"route" binding:"required"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
