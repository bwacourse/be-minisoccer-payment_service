package repositories

import (
	paymentRepository "payment-service/repositories/payment"
	paymentHistoryRepository "payment-service/repositories/paymenthistory"

	"gorm.io/gorm"
)

type RegistryRepository struct {
	db *gorm.DB
}

type IRegistryRepository interface {
	GetPayment() paymentRepository.IPaymentRepository
	GetPaymentHistory() paymentHistoryRepository.IPaymentHistoryRepository
	GetTx() *gorm.DB
}

func NewRegistryRepository(db *gorm.DB) IRegistryRepository {
	return &RegistryRepository{
		db: db,
	}
}

// GetPayment implements [IRegistryRepository].
func (r *RegistryRepository) GetPayment() paymentRepository.IPaymentRepository {
	return paymentRepository.NewPaymentRepository(r.db)
}

// GetPaymentHistory implements [IRegistryRepository].
func (r *RegistryRepository) GetPaymentHistory() paymentHistoryRepository.IPaymentHistoryRepository {
	return paymentHistoryRepository.NewPaymentHistoryRepository(r.db)
}

// GetTx implements [IRegistryRepository].
func (r *RegistryRepository) GetTx() *gorm.DB {
	return r.db
}
