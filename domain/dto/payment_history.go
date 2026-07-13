package dto

import "payment-service/constants"

type PaymentHistoryRequest struct {
	PaymentID uint                          `json:"paymentId"`
	Status    constants.PaymentStatusString `json:"status"`
}
