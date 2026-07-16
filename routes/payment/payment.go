package routes

import (
	"payment-service/clients"
	"payment-service/constants"
	controllers "payment-service/controllers/http"
	"payment-service/middlewares"

	"github.com/gin-gonic/gin"
)

type PaymentRoute struct {
	controller  controllers.IRegistryController
	client      clients.IClientRegistry
	routerGroup *gin.RouterGroup
}

type IPaymentRoute interface {
	Run()
}

func NewPaymentRoute(controller controllers.IRegistryController, client clients.IClientRegistry, routerGroup *gin.RouterGroup) IPaymentRoute {
	return &PaymentRoute{
		controller:  controller,
		client:      client,
		routerGroup: routerGroup,
	}
}

func (r *PaymentRoute) Run() {
	group := r.routerGroup.Group("/payment")

	group.POST("/webhook", r.controller.GetPayment().Webhook)

	group.Use(middlewares.Authenticate())

	group.GET("", middlewares.CheckRole([]string{
		constants.Admin,
		constants.Customer,
	}, r.client), r.controller.GetPayment().GetAllWithPagination)

	group.GET("/:uuid", middlewares.CheckRole([]string{
		constants.Admin,
		constants.Customer,
	}, r.client), r.controller.GetPayment().GetByUUID)

	group.POST("", middlewares.CheckRole([]string{
		constants.Customer,
	}, r.client), r.controller.GetPayment().Create)
}
