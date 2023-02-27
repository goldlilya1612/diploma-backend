package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Lecturer struct {
	UserID uuid.UUID      `gorm:"type:uuid;uniqueIndex"`
	Name   string         `gorm:"type:varchar(255);not null"`
	Groups pq.StringArray `gorm:"type:varchar(255);not null"`

	User User
}
