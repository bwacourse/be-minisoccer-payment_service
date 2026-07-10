package repositories

import "gorm.io/gorm"

type RegistryRepository struct {
	db *gorm.DB
}

type IRegistryRepository interface {
	GetPayment() any
}

func NewRegistryRepository(db *gorm.DB) IRegistryRepository {
	return &RegistryRepository{
		db: db,
	}
}

// GetField implements IRegistryRepository.
func (r *RegistryRepository) GetPayment() any {
	panic("not implemented")
}
