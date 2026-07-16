package controllers

import (
	"net/http"
	"payment-service/common/response"
	"payment-service/domain/dto"
	"payment-service/services"

	errValidation "payment-service/common/error"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type PaymentController struct {
	service services.IRegistryService
}

type IPaymentController interface {
	GetAllWithPagination(*gin.Context)
	GetByUUID(*gin.Context)
	Create(*gin.Context)
	Webhook(*gin.Context)
}

func NewPaymentController(service services.IRegistryService) IPaymentController {
	return &PaymentController{service: service}
}

// GetAllWithPagination implements [IPaymentController].
func (p *PaymentController) GetAllWithPagination(ctx *gin.Context) {
	var params dto.PaymentRequestParam

	err := ctx.ShouldBindQuery(&params)

	if err != nil {
		response.HttpResponse(response.ParamHTTPResponse{
			Gin:  ctx,
			Code: http.StatusBadRequest,
			Err:  err,
		})
		return
	}

	validate := validator.New()

	// Validate the request parameters
	if err = validate.Struct(params); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)

		response.HttpResponse(response.ParamHTTPResponse{
			Code:    http.StatusBadRequest,
			Err:     err,
			Message: &errMessage,
			Data:    errResponse,
			Gin:     ctx,
		})
		return
	}

	result, err := p.service.GetPayment().GetAllWithPagination(ctx, &params)

	if err != nil {
		response.HttpResponse(response.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: result,
		Gin:  ctx,
	})
}

// GetByUUID implements [IPaymentController].
func (p *PaymentController) GetByUUID(ctx *gin.Context) {
	result, err := p.service.GetPayment().GetByUUID(ctx, ctx.Param("uuid"))

	if err != nil {
		response.HttpResponse(response.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: result,
		Gin:  ctx,
	})
}

// Create implements [IPaymentController].
func (p *PaymentController) Create(ctx *gin.Context) {
	var request dto.PaymentRequest

	err := ctx.ShouldBindJSON(&request)

	if err != nil {
		response.HttpResponse(response.ParamHTTPResponse{
			Gin:  ctx,
			Code: http.StatusBadRequest,
			Err:  err,
		})
		return
	}

	validate := validator.New()

	// Validate the request parameters
	if err = validate.Struct(request); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)

		response.HttpResponse(response.ParamHTTPResponse{
			Code:    http.StatusBadRequest,
			Err:     err,
			Message: &errMessage,
			Data:    errResponse,
			Gin:     ctx,
		})
		return
	}

	result, err := p.service.GetPayment().Create(ctx, &request)

	if err != nil {
		response.HttpResponse(response.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResponse{
		Code: http.StatusCreated,
		Data: result,
		Gin:  ctx,
	})
}

// Webhook implements [IPaymentController].
func (p *PaymentController) Webhook(ctx *gin.Context) {
	var request dto.WebHook

	err := ctx.ShouldBindJSON(&request)

	if err != nil {
		response.HttpResponse(response.ParamHTTPResponse{
			Gin:  ctx,
			Code: http.StatusBadRequest,
			Err:  err,
		})
		return
	}

	err = p.service.GetPayment().Webhook(ctx, &request)

	if err != nil {
		response.HttpResponse(response.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Gin:  ctx,
	})
}
