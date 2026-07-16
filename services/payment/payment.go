package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	clients "payment-service/clients/midtrans"
	"payment-service/common/gcs"
	"payment-service/common/util"
	"payment-service/config"
	"payment-service/constants"
	"payment-service/controllers/kafka"
	"payment-service/domain/dto"
	"payment-service/domain/models"
	"payment-service/repositories"
	"strings"
	"time"

	"gorm.io/gorm"

	// errWrap "payment-service/common/error"
	// errConstant "payment-service/constants/error"
	errPayment "payment-service/constants/error/payment"
)

type PaymentService struct {
	repository repositories.IRegistryRepository
	gcs        gcs.IGCSClient
	kafka      kafka.IKafkaRegistry
	midtrans   clients.IMidtransClient
}

type IPaymentService interface {
	GetAllWithPagination(context.Context, *dto.PaymentRequestParam) (*util.PaginationResult, error)
	GetByUUID(context.Context, string) (*dto.PaymentResponse, error)
	Create(context.Context, *dto.PaymentRequest) (*dto.PaymentResponse, error)
	Webhook(context.Context, *dto.WebHook) error
}

func NewPaymentService(
	repository repositories.IRegistryRepository,
	gcs gcs.IGCSClient,
	kafka kafka.IKafkaRegistry,
	midtrans clients.IMidtransClient,
) IPaymentService {
	return &PaymentService{
		repository: repository,
		gcs:        gcs,
		kafka:      kafka,
		midtrans:   midtrans,
	}
}

// GetAllWithPagination implements [IPaymentService].
func (p *PaymentService) GetAllWithPagination(ctx context.Context, params *dto.PaymentRequestParam) (*util.PaginationResult, error) {
	payments, total, err := p.repository.GetPayment().FindAllWithPagination(ctx, params)

	if err != nil {
		return nil, err
	}

	var paymentResponses = make([]dto.PaymentResponse, 0, len(payments))

	for _, payment := range payments {
		paymentResponses = append(paymentResponses, dto.PaymentResponse{
			UUID:          payment.UUID,
			TransactionID: payment.TransactionID,
			OrderID:       payment.OrderID,
			Amount:        payment.Amount,
			Status:        payment.Status.GetStatusString(),
			PaymentLink:   payment.PaymentLink,
			InvoiceLink:   payment.InvoiceLink,
			VANumber:      payment.VANumber,
			Bank:          payment.Bank,
			Description:   payment.Description,
			ExpiredAt:     payment.ExpiredAt,
			CreatedAt:     payment.CreatedAt,
			UpdatedAt:     payment.UpdatedAt,
		})
	}

	paginationParam := util.PaginationParam{
		Page:  params.Page,
		Limit: params.Limit,
		Count: total,
		Data:  paymentResponses,
	}

	response := util.GeneratePagination(paginationParam)

	return &response, nil
}

// GetByUUID implements [IPaymentService].
func (p *PaymentService) GetByUUID(ctx context.Context, uuid string) (*dto.PaymentResponse, error) {
	payment, err := p.repository.GetPayment().FindByUUID(ctx, uuid)

	if err != nil {
		return nil, err
	}

	return &dto.PaymentResponse{
		UUID:          payment.UUID,
		TransactionID: payment.TransactionID,
		OrderID:       payment.OrderID,
		Amount:        payment.Amount,
		Status:        payment.Status.GetStatusString(),
		PaymentLink:   payment.PaymentLink,
		InvoiceLink:   payment.InvoiceLink,
		VANumber:      payment.VANumber,
		Bank:          payment.Bank,
		Description:   payment.Description,
		ExpiredAt:     payment.ExpiredAt,
		CreatedAt:     payment.CreatedAt,
		UpdatedAt:     payment.UpdatedAt,
	}, nil
}

// Create implements [IPaymentService].
func (p *PaymentService) Create(ctx context.Context, request *dto.PaymentRequest) (*dto.PaymentResponse, error) {
	var (
		txErr, err error
		payment    *models.Payment
		response   *dto.PaymentResponse
		midtrans   *clients.MidtransData
	)

	err = p.repository.GetTx().Transaction(func(tx *gorm.DB) error {
		if !request.ExpiredAt.After(time.Now()) {
			return errPayment.ErrExpireAtInvalid
		}

		midtrans, txErr = p.midtrans.CreatePaymentLink(request)

		if txErr != nil {
			return txErr
		}

		paymentRequest := &dto.PaymentRequest{
			OrderID:     request.OrderID,
			Amount:      request.Amount,
			Description: request.Description,
			ExpiredAt:   request.ExpiredAt,
			PaymentLink: midtrans.RedirectURL,
		}

		payment, txErr = p.repository.GetPayment().Create(ctx, tx, paymentRequest)

		return nil
	})

	if err != nil {
		return nil, err
	}

	response = &dto.PaymentResponse{
		UUID:        payment.UUID,
		OrderID:     payment.OrderID,
		Amount:      payment.Amount,
		Status:      payment.Status.GetStatusString(),
		PaymentLink: payment.PaymentLink,
		Description: payment.Description,
	}

	return response, txErr
}

// Webhook implements [IPaymentService].
func (p *PaymentService) Webhook(ctx context.Context, request *dto.WebHook) error {
	var (
		txErr, err         error
		paymentAfterUpdate *models.Payment
		paidAt             *time.Time
		invoiceLink        string
		pdf                []byte
	)

	// Begin Transaction
	err = p.repository.GetTx().Transaction(func(tx *gorm.DB) error {
		// Check payment
		_, txErr = p.repository.GetPayment().FindByOrderID(ctx, request.OrderID.String())

		if txErr != nil {
			return txErr
		}

		// Set paidAt if transaction status is settlement
		if request.TransactionStatus == constants.SettlementString {
			now := time.Now()
			paidAt = &now
		}

		status := request.TransactionStatus.GetStatusInt()

		vaNumber := request.VANumbers[0].VANumber
		bank := request.VANumbers[0].Bank

		// Update Payment
		_, txErr = p.repository.GetPayment().Update(ctx, tx, request.OrderID.String(), &dto.UpdatePaymentRequest{
			TransactionID: &request.TransactionID,
			Status:        &status,
			PaidAt:        paidAt,
			VANumber:      &vaNumber,
			Bank:          &bank,
			Acquirer:      request.Acquirer,
		})

		if txErr != nil {
			return txErr
		}

		// Get detail payment by orderId which has been updated
		paymentAfterUpdate, txErr = p.repository.GetPayment().FindByOrderID(ctx, request.OrderID.String())

		if txErr != nil {
			return txErr
		}

		// Create Payment History
		txErr = p.repository.GetPaymentHistory().Create(ctx, tx, &dto.PaymentHistoryRequest{
			PaymentID: paymentAfterUpdate.ID,
			Status:    paymentAfterUpdate.Status.GetStatusString(),
		})

		// Generate invoice if transaction status is settlement
		if request.TransactionStatus == constants.SettlementString {
			paidDay := paidAt.Format("02")
			paidMonth := p.convertToIndonesianMonth(paidAt.Format("January"))
			paidYear := paidAt.Format("2006")

			invoiceNumber := fmt.Sprintf("INV/%s/ORD/%d", time.Now().Format(time.DateOnly), p.randomNumber())

			total := util.RupiahFormat(&paymentAfterUpdate.Amount)

			// Invoice Request
			invoiceRequest := &dto.InvoiceRequest{
				InvoiceNumber: invoiceNumber,
				Data: dto.InvoiceData{
					PaymentDetail: dto.InvoicePaymentDetail{
						PaymentMethod: request.PaymentType,
						BankName:      strings.ToUpper(*paymentAfterUpdate.Bank),
						VANumber:      *paymentAfterUpdate.VANumber,
						Date:          fmt.Sprintf("%s %s %s", paidDay, paidMonth, paidYear),
						IsPaid:        true,
					},
					Items: []dto.InvoiceItem{
						{
							Description: *paymentAfterUpdate.Description,
							Price:       total,
						},
					},
					Total: total,
				},
			}

			// Generate PDF
			if pdf, txErr = p.generatePDF(invoiceRequest); txErr != nil {
				return txErr
			}

			// Upload to GCS
			if invoiceLink, txErr = p.uploadToGCS(ctx, invoiceNumber, pdf); txErr != nil {
				return txErr
			}

			// Update Payment
			paymentAfterUpdate, txErr = p.repository.GetPayment().Update(ctx, tx, request.OrderID.String(), &dto.UpdatePaymentRequest{
				InvoiceLink: &invoiceLink,
			})

			if txErr != nil {
				return txErr
			}

			// TODO: Send Email
		}

		return nil
	}) // End Transaction process

	if err != nil {
		return err
	}

	// Produce to Kafka
	if err = p.produceToKafka(request, paymentAfterUpdate, paidAt); err != nil {
		return err
	}

	return nil
}

/*
*
** Private Methods **
*
 */

// Convert English month to Indonesian month
func (p *PaymentService) convertToIndonesianMonth(englishMonth string) string {
	monthMap := map[string]string{
		"January":   "Januari",
		"February":  "Februari",
		"March":     "Maret",
		"April":     "April",
		"May":       "Mei",
		"June":      "Juni",
		"July":      "Juli",
		"August":    "Agustus",
		"September": "September",
		"October":   "Oktober",
		"November":  "November",
		"December":  "Desember",
	}

	indonesianMonth, ok := monthMap[englishMonth]

	if !ok {
		return errors.New("month not found").Error()
	}

	return indonesianMonth
}

// Generate PDF from HTML template
func (p *PaymentService) generatePDF(request *dto.InvoiceRequest) ([]byte, error) {
	// Template file name
	htmlTemplatePath := "template/invoice.html"

	// Read html template
	htmlTemplate, err := os.ReadFile(htmlTemplatePath)

	if err != nil {
		return nil, err
	}

	var data map[string]interface{}

	// Convert struct to map[string]interface{}
	jsonData, _ := json.Marshal(request)

	// Unmarshal map[string]interface{} to data
	err = json.Unmarshal(jsonData, &data)

	if err != nil {
		return nil, err
	}

	// Generate PDF from HTML template
	pdf, err := util.GeneratePDFFromHTML(string(htmlTemplate), data)

	return pdf, nil
}

// Upload PDF to GCS
func (p *PaymentService) uploadToGCS(ctx context.Context, invoiceNumber string, pdf []byte) (string, error) {
	invoiceNumberReplace := strings.ToLower(strings.ReplaceAll(invoiceNumber, "/", "-"))

	filename := fmt.Sprintf("%s.pdf", invoiceNumberReplace)

	url, err := p.gcs.UploadFile(ctx, filename, pdf)

	if err != nil {
		return "", err
	}

	return url, nil
}

// Generate random number for invoice
func (p *PaymentService) randomNumber() int {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	number := random.Intn(900000) + 100000

	return number
}

// Map transaction status to event
func (p *PaymentService) mapTransactionStatusToEvent(status constants.PaymentStatusString) string {
	var paymentStatus string

	switch status {
	case constants.PendingString:
		paymentStatus = strings.ToUpper(constants.PendingString.String())
	case constants.SettlementString:
		paymentStatus = strings.ToUpper(constants.SettlementString.String())
	case constants.ExpireString:
		paymentStatus = strings.ToUpper(constants.ExpireString.String())
	}

	return paymentStatus
}

// Produce message to Kafka
func (p *PaymentService) produceToKafka(
	req *dto.WebHook,
	payment *models.Payment,
	paidAt *time.Time,
) error {
	event := dto.KafkaEvent{
		Name: p.mapTransactionStatusToEvent(req.TransactionStatus),
	}

	metadata := dto.KafkaMetaData{
		Sender:    "payment-service",
		SendingAt: time.Now().Format(time.RFC3339),
	}

	body := dto.KafkaBody{
		Type: "JSON",
		Data: &dto.KafkaData{
			OrderID:   payment.OrderID,
			PaymentID: payment.UUID,
			Status:    req.TransactionStatus.String(),
			PaidAt:    paidAt,
			ExpiredAt: *payment.ExpiredAt,
		},
	}

	kafkaMessage := dto.KafkaMessage{
		Event:    event,
		Metadata: metadata,
		Body:     body,
	}

	topic := config.Config.Kafka.Topic

	kafkaMessageJSON, err := json.Marshal(kafkaMessage)

	if err != nil {
		return err
	}

	err = p.kafka.GetKafkaProducer().ProduceMessage(topic, kafkaMessageJSON)

	if err != nil {
		return err
	}

	return nil
}
