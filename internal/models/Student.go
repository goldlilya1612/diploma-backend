package models

import (
	"github.com/google/uuid"
)

type Student struct {
	UserID uuid.UUID `gorm:"type:uuid;uniqueIndex"`
	Name   string    `gorm:"type:varchar(255);not null"`
	Group  string    `gorm:"type:varchar(255);not null"`

	User User
}
