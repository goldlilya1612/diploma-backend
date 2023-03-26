package models

type Article struct {
	ID        int    `gorm:"primaryKey;uniqueIndex"`
	Name      string `gorm:"type:varchar(255);not null"`
	ChapterID int    `gorm:"type:integer;not null"`

	Chapter Chapter `gorm:"foreignKey:ChapterID"`
}
