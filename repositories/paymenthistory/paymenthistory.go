package repositories

import (
	"context"
	"payment-service/domain/dto"
	"payment-service/domain/models"

	"gorm.io/gorm"

	errWrap "payment-service/common/error"
	errConstant "payment-service/constants/error"
)

type PaymentHistory struct {
	db *gorm.DB
}

type IPaymentHistoryRepository interface {
	Create(context.Context, *gorm.DB, *dto.PaymentHistoryRequest) error
}

func NewPaymentHistoryRepository(db *gorm.DB) IPaymentHistoryRepository {
	return &PaymentHistory{
		db: db,
	}
}

// Create implements [IPaymentHistoryRepository].
func (p *PaymentHistory) Create(ctx context.Context, tx *gorm.DB, request *dto.PaymentHistoryRequest) error {
	paymentHistory := models.PaymentHistory{
		PaymentID: request.PaymentID,
		Status:    request.Status,
	}

	if tx == nil {
		tx = p.db
	}

	err := tx.WithContext(ctx).Create(&paymentHistory).Error

	if err != nil {
		return errWrap.WrapError(errConstant.ErrSQLError)
	}

	return nil
}
