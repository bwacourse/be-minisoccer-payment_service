package controllers

type RegistryController struct {
	service any
}

type IRegistryController interface {
	GetPayment() any
}

func NewRegistryController(service any) IRegistryController {
	return &RegistryController{
		service: service,
	}
}

// GetField implements IRegistryController.
func (c *RegistryController) GetPayment() any {
	panic("not implemented")
}
