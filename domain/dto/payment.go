package dto

import (
	"time"

	"github.com/google/uuid"
)

type PaymentRequest struct {
	Status string `json:"status" validate:"required"`
}

type PaymentResponse struct {
	UUID      uuid.UUID  `json:"uuid"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}
