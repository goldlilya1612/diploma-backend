package models

type Chapter struct {
	ID       int    `gorm:"primaryKey;uniqueIndex"`
	Name     string `gorm:"type:varchar(255);not null"`
	CourseID int    `gorm:"type:integer;not null"`

	Course Course `gorm:"foreignKey:CourseID"`
}
