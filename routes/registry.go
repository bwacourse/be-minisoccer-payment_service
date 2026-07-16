package routes

import (
	"payment-service/clients"
	controllers "payment-service/controllers/http"
	routes "payment-service/routes/payment"

	"github.com/gin-gonic/gin"
)

type RegistryRoute struct {
	controllers controllers.IRegistryController
	routerGroup *gin.RouterGroup
	client      clients.IClientRegistry
}

type IRegistryRoute interface {
	Serve()
}

func NewRegistryRoute(controller controllers.IRegistryController, routerGroup *gin.RouterGroup, client clients.IClientRegistry) IRegistryRoute {
	return &RegistryRoute{
		controllers: controller,
		routerGroup: routerGroup,
		client:      client,
	}
}

func (c *RegistryRoute) Serve() {
	c.paymentRoute().Run()
}

func (c *RegistryRoute) paymentRoute() routes.IPaymentRoute {
	return routes.NewPaymentRoute(c.controllers, c.client, c.routerGroup)
}
