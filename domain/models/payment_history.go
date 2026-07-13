package models

import (
	"payment-service/constants"
)

type PaymentHistory struct {
	ID        uint                          `gorm:"primaryKey;autoIncrement"`
	PaymentID uint                          `gorm:"type:bigint;not null"`
	Status    constants.PaymentStatusString `gorm:"type:varchar(50);not null"`
	CreatedAt string                        `gorm:"type:timestamp;not null"`
	UpdatedAt string                        `gorm:"type:timestamp;not null"`
}
