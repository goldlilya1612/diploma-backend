package models

type Image struct {
	ID   int    `gorm:"primaryKey;uniqueIndex"`
	Name string `gorm:"type:varchar(255);not null"`
	Path string `gorm:"type:varchar(255);not null"`
	URL  string `gorm:"type:text;not null"`
}
