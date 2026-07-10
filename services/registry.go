package services

type RegistryService struct {
	gcs        any
	repository any
}

type IRegistryService interface {
	GetPayment() any
}

func NewRegistryService(repository any, gcs any) IRegistryService {
	return &RegistryService{
		repository: repository,
		gcs:        gcs,
	}
}

// GetField implements Service.
func (s *RegistryService) GetPayment() any {
	panic("not implemented")
}
