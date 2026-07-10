package models

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UUID      uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
