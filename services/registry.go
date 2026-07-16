package services

import (
	clients "payment-service/clients/midtrans"
	"payment-service/common/gcs"
	"payment-service/controllers/kafka"
	"payment-service/repositories"
	services "payment-service/services/payment"
)

type RegistryService struct {
	repository repositories.IRegistryRepository
	gcs        gcs.IGCSClient
	kafka      kafka.IKafkaRegistry
	midtrans   clients.IMidtransClient
}

type IRegistryService interface {
	GetPayment() services.IPaymentService
}

func NewRegistryService(
	repository repositories.IRegistryRepository,
	gcs gcs.IGCSClient,
	kafka kafka.IKafkaRegistry,
	midtrans clients.IMidtransClient,
) IRegistryService {
	return &RegistryService{
		repository: repository,
		gcs:        gcs,
		kafka:      kafka,
		midtrans:   midtrans,
	}
}

// GetPayment implements [IRegistryService].
func (r *RegistryService) GetPayment() services.IPaymentService {
	return services.NewPaymentService(r.repository, r.gcs, r.kafka, r.midtrans)
}
