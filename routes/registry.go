package routes

import (
	"payment-service/clients"
	"payment-service/controllers"

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
	panic("not implemented")
}
