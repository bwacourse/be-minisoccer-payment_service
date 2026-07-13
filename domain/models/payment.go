package models

import (
	"payment-service/constants"
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID               uint                     `gorm:"primaryKey;autoIncrement"`
	UUID             uuid.UUID                `gorm:"type:uuid;not null"`
	Amount           float64                  `gorm:"not null"`
	Status           *constants.PaymentStatus `gorm:"not null"`
	PaymentLink      string                   `gorm:"type:varchar(255);not null"`
	InvoiceLink      *string                  `gorm:"type:varchar(255);default;null"`
	VANumber         *string                  `gorm:"type:varchar(60);default;null"`
	Bank             *string                  `gorm:"type:varchar(100);default;null"`
	Acquirer         *string                  `gorm:"type:varchar(100);default;null"`
	TransactionID    *string                  `gorm:"type:varchar(100);default;null"`
	Description      *string                  `gorm:"type:text;default;null"`
	PaidAt           *time.Time               `gorm:"type:timestamp"`
	ExpiredAt        time.Time                `gorm:"type:timestamp"`
	CreatedAt        time.Time                `gorm:"type:timestamp;not null"`
	UpdatedAt        time.Time                `gorm:"type:timestamp;not null"`
	PaymentHistories []PaymentHistory         `gorm:"foreignKey:payment_id;references:id;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
}
