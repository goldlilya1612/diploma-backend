package models

import "time"

type Article struct {
	ID        int    `gorm:"primaryKey;uniqueIndex"`
	Name      string `gorm:"type:varchar(255);not null"`
	ChapterID int    `gorm:"type:integer;not null"`

	Chapter Chapter `gorm:"foreignKey:ChapterID"`
}

type ArticleResponse struct {
	ID        int    `json:"id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	ChapterID int    `json:"courseID" binding:"required"`
	Route     string `json:"route" binding:"required"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
