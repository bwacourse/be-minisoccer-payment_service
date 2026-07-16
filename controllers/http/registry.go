package controllers

import (
	controllers "payment-service/controllers/http/payment"
	"payment-service/services"
)

type RegistryController struct {
	service services.IRegistryService
}

type IRegistryController interface {
	GetPayment() controllers.IPaymentController
}

func NewRegistryController(service services.IRegistryService) IRegistryController {
	return &RegistryController{service: service}
}

// GetPayment implements [IRegistryController]
func (r *RegistryController) GetPayment() controllers.IPaymentController {
	return controllers.NewPaymentController(r.service)
}
