package models

type Image struct {
	ID   int    `json:"id" gorm:"primaryKey;uniqueIndex"`
	Name string `json:"name" gorm:"type:varchar(255);not null"`
	Path string `json:"path" gorm:"type:varchar(255);not null"`
	URL  string `json:"url" gorm:"type:text;not null"`
}
